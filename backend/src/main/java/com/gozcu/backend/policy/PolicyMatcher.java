package com.gozcu.backend.policy;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import java.util.List;
import java.util.regex.Pattern;

@Service
@RequiredArgsConstructor
public class PolicyMatcher {

    private final CommandPolicyRepository repository;

    /**
     * Hibrit eşleştirme: önce tam eşleşme, yoksa wildcard fallback.
     * Hiçbir policy eşleşmezse -> LOW (otomatik geçer, onay istemez).
     */
    public RiskLevel resolveRisk(String command) {
        List<CommandPolicy> policies = repository.findAll();

        for (CommandPolicy p : policies) {
            if (!p.getPattern().contains("*") && p.getPattern().equals(command)) {
                return p.getRiskLevel();
            }
        }

        for (CommandPolicy p : policies) {
            if (p.getPattern().contains("*") && toRegex(p.getPattern()).matcher(command).matches()) {
                return p.getRiskLevel();
            }
        }

        return RiskLevel.LOW;
    }

    private Pattern toRegex(String wildcardPattern) {
        String[] parts = wildcardPattern.split("\\*");
        StringBuilder regex = new StringBuilder();
        for (int i = 0; i < parts.length; i++) {
            regex.append(Pattern.quote(parts[i]));
            if (i < parts.length - 1) {
                regex.append(".*");
            }
        }
        if (wildcardPattern.endsWith("*")) {
            regex.append(".*");
        }
        return Pattern.compile("^" + regex + "$");
    }
}