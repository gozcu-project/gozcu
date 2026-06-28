package com.gozcu.backend.policy;

import com.gozcu.backend.policy.CommandPolicy;

import org.springframework.data.jpa.repository.JpaRepository;
import java.util.List;

public interface CommandPolicyRepository extends JpaRepository<CommandPolicy, Long> {
    List<CommandPolicy> findAllByPatternContaining(String wildcard);
}