package com.gozcu.backend.policy.dto;

import com.gozcu.backend.policy.RiskLevel;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class PolicyCreateDto {

    @NotBlank
    private String pattern;

    @NotNull
    private RiskLevel riskLevel;
}