package sieve

import (
	"regexp"
	"sync"
)

// Risk is the classification returned for a scanned input.
type Risk string

const (
	RiskSafe       Risk = "safe"
	RiskSuspicious Risk = "suspicious"
	RiskBlocked    Risk = "blocked"
)

// Match describes a single pattern that fired during a scan.
type Match struct {
	Category string `json:"category"`
	Pattern  string `json:"pattern"`
	Risk     Risk   `json:"risk"`
}

// Result is the outcome of scanning an input string.
type Result struct {
	Risk    Risk    `json:"risk"`
	Matches []Match `json:"matches,omitempty"`
	Input   string  `json:"input,omitempty"`
}

// Rule is a single detection rule.
type Rule struct {
	Category string
	Pattern  string
	Risk     Risk
	re       *regexp.Regexp
}

// FlagLogger receives structured entries for non-safe scan results.
// Implement in-process (e.g. flaglog) so library users can log without the HTTP sidecar.
type FlagLogger interface {
	Write(result Result, source string) error
}

// Config controls scan behavior.
type Config struct {
	BlockOn        Risk       // inputs at or above this risk are blocked (default: blocked)
	IncludeInput   bool       // echo the input back in Result
	CustomRules    []Rule     // appended to the built-in library
	DisableBuiltin bool
	Logger         FlagLogger // optional; Write is called for suspicious/blocked results
}

// Sieve screens agent inputs for prompt injection patterns.
type Sieve struct {
	mu    sync.RWMutex
	rules []Rule
	cfg   Config
}

// New builds a Sieve with built-in rules plus any custom rules from cfg.
func New(cfg Config) (*Sieve, error) {
	if cfg.BlockOn == "" {
		cfg.BlockOn = RiskBlocked
	}
	s := &Sieve{cfg: cfg}
	var rules []Rule
	if !cfg.DisableBuiltin {
		rules = append(rules, builtinRules()...)
	}
	rules = append(rules, cfg.CustomRules...)
	compiled, err := compileRules(rules)
	if err != nil {
		return nil, err
	}
	s.rules = compiled
	return s, nil
}

// Scan classifies input and returns a risk result.
func (s *Sieve) Scan(input string) Result {
	return s.ScanWithSource(input, "")
}

// ScanWithSource is Scan with an optional source label for FlagLogger.
func (s *Sieve) ScanWithSource(input, source string) Result {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalized := normalize(input)
	var matches []Match
	highest := RiskSafe

	for _, rule := range s.rules {
		if rule.re.MatchString(normalized) {
			matches = append(matches, Match{
				Category: rule.Category,
				Pattern:  rule.Pattern,
				Risk:     rule.Risk,
			})
			if riskRank(rule.Risk) > riskRank(highest) {
				highest = rule.Risk
			}
		}
	}

	out := Result{Risk: highest, Matches: matches}
	if s.cfg.IncludeInput {
		out.Input = input
	}
	if s.cfg.Logger != nil && out.Risk != RiskSafe {
		_ = s.cfg.Logger.Write(out, source)
	}
	return out
}

// ShouldBlock reports whether the result should be rejected before the model call.
func (s *Sieve) ShouldBlock(r Result) bool {
	return riskRank(r.Risk) >= riskRank(s.cfg.BlockOn)
}

func riskRank(r Risk) int {
	switch r {
	case RiskBlocked:
		return 2
	case RiskSuspicious:
		return 1
	default:
		return 0
	}
}

func compileRules(rules []Rule) ([]Rule, error) {
	out := make([]Rule, 0, len(rules))
	for _, rule := range rules {
		re, err := regexp.Compile("(?i)" + rule.Pattern)
		if err != nil {
			return nil, err
		}
		rule.re = re
		if rule.Risk == "" {
			rule.Risk = RiskSuspicious
		}
		out = append(out, rule)
	}
	return out, nil
}

func builtinRules() []Rule {
	return []Rule{
		{Category: "role_override", Pattern: `\b(you are now|act as|pretend (to be|you('re| are))|assume the role)\b`, Risk: RiskBlocked},
		{Category: "ignore_instructions", Pattern: `\b(ignore|disregard|forget)\b.{0,80}\b(previous|prior|above|system)\b.{0,40}\b(instructions?|prompts?|rules?)\b`, Risk: RiskBlocked},
		{Category: "jailbreak", Pattern: `\b(jailbreak|dan mode|developer mode|do anything now)\b`, Risk: RiskBlocked},
		{Category: "system_prompt_leak", Pattern: `\b(reveal|show|print|dump)\b.{0,30}\b(system prompt|hidden (instructions?|prompt)|initial instructions?)\b`, Risk: RiskBlocked},
		{Category: "delimiter_injection", Pattern: "(?s)(" + "```" + `|\{\{|\[INST\]|<\|im_start\|>|<system>)`, Risk: RiskSuspicious},
		{Category: "policy_bypass", Pattern: `\b(bypass|circumvent|override)\b.{0,30}\b(safety|filter|guardrail|policy|restriction)\b`, Risk: RiskBlocked},
		{Category: "encoding_obfuscation", Pattern: `\b(base64|rot13|hex encoded)\b.{0,20}\b(decode|decode this|interpret)\b`, Risk: RiskSuspicious},
		{Category: "tool_abuse", Pattern: `\b(exfiltrat|steal|exfil)\b.{0,40}\b(api key|secret|credential|token|password)\b`, Risk: RiskBlocked},
		{Category: "contextual_persuasion", Pattern: `\b(this is (just |only )?(a )?(test|hypothetical|roleplay|fiction))\b.{0,40}\b(ignore|bypass|disable)\b`, Risk: RiskSuspicious},
		{Category: "instruction_smuggle", Pattern: `(?s)<\s*(system|assistant|instruction)\s*>`, Risk: RiskBlocked},
		{Category: "multiline_ignore", Pattern: `(?s)\b(ignore|disregard)\b.{0,120}\b(all|any)\b.{0,40}\b(rules?|guidelines?|safeguards?)\b`, Risk: RiskBlocked},
		{Category: "translated_jailbreak", Pattern: `\b(ignora|ignorez|ignoriere|ignora le|desconsidera)\b.{0,40}\b(instrucciones|instructions|anweisungen|instruções|regras)\b`, Risk: RiskBlocked},
		{Category: "homoglyph_hint", Pattern: `(?i)(іgnore|іnstructions|ѕystem prompt|јailbreak)`, Risk: RiskSuspicious},
		{Category: "hidden_html_instruction", Pattern: `(?s)(<!--.{0,200}\b(ignore|system prompt|you are now|jailbreak)\b.{0,200}-->|style\s*=\s*["'][^"']*display\s*:\s*none)`, Risk: RiskSuspicious},
		{Category: "sudo_mode", Pattern: `\b(sudo mode|admin mode|god mode|unrestricted mode)\b`, Risk: RiskBlocked},
		{Category: "end_of_prompt", Pattern: `\b(end of (system )?prompt|new (system )?instructions?:)\b`, Risk: RiskBlocked},
	}
}

// BuiltinCategoryCount is the number of distinct built-in rule categories.
func BuiltinCategoryCount() int {
	seen := map[string]struct{}{}
	for _, r := range builtinRules() {
		seen[r.Category] = struct{}{}
	}
	return len(seen)
}
