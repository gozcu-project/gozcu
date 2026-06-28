package com.gozcu.backend.approval.dto;

import com.gozcu.backend.approval.ApprovalStatus;
import com.gozcu.backend.policy.RiskLevel;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;

import java.time.Instant;

@Getter
@Builder
@AllArgsConstructor
public class ApprovalResponseDto {
    private Long id;
    private String requestedBy;
    private String hostName;
    private String command;
    private RiskLevel riskLevel;
    private ApprovalStatus status;
    private Instant createdAt;
    private Instant resolvedAt;
    private String resolvedBy;
}