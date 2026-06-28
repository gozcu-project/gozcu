package com.gozcu.backend.policy;

import com.gozcu.backend.policy.dto.PolicyCreateDto;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import java.util.List;

@RestController
@RequestMapping("/api/policies")
@RequiredArgsConstructor
public class PolicyController {

    private final CommandPolicyRepository repository;

    @PostMapping
    public ResponseEntity<CommandPolicy> create(@RequestBody @Valid PolicyCreateDto dto) {
        CommandPolicy policy = CommandPolicy.builder()
                .pattern(dto.getPattern())
                .riskLevel(dto.getRiskLevel())
                .build();
        return ResponseEntity.ok(repository.save(policy));
    }

    @GetMapping
    public List<CommandPolicy> list() {
        return repository.findAll();
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> delete(@PathVariable Long id) {
        repository.deleteById(id);
        return ResponseEntity.noContent().build();
    }
}