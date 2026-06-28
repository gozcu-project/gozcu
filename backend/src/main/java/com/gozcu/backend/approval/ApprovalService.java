package com.gozcu.backend.approval;

import com.gozcu.backend.approval.dto.ApprovalCreateRequestDto;
import com.gozcu.backend.approval.dto.ApprovalResponseDto;
import com.gozcu.backend.audit.AuditService;
import com.gozcu.backend.policy.PolicyMatcher;
import com.gozcu.backend.policy.RiskLevel;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.context.request.async.DeferredResult;

import java.time.Instant;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

@Service
@RequiredArgsConstructor
public class ApprovalService {

    private static final long WAIT_TIMEOUT_MS = 30_000L;
    private static final String AUTO_RESOLVER = "system:auto-low-risk";

    private final ApprovalRequestRepository repository;
    private final AuditService auditService;
    private final PolicyMatcher policyMatcher;

    private final Map<Long, DeferredResult<ResponseEntity<ApprovalResponseDto>>> waiters =
            new ConcurrentHashMap<>();

    public ApprovalResponseDto create(ApprovalCreateRequestDto dto) {
        RiskLevel risk = policyMatcher.resolveRisk(dto.getCommand());

        ApprovalRequest.ApprovalRequestBuilder builder = ApprovalRequest.builder()
                .requestedBy(dto.getRequestedBy())
                .hostName(dto.getHostName())
                .command(dto.getCommand())
                .riskLevel(risk)
                .createdAt(Instant.now());

        ApprovalRequest req;
        if (risk == RiskLevel.LOW) {
            // Eşleşmeyen veya LOW olarak tanımlı komutlar onay istemeden geçer
            req = builder
                    .status(ApprovalStatus.APPROVED)
                    .resolvedAt(Instant.now())
                    .resolvedBy(AUTO_RESOLVER)
                    .build();
        } else {
            req = builder
                    .status(ApprovalStatus.PENDING)
                    .build();
        }

        repository.save(req);

        String action = (risk == RiskLevel.LOW) ? "AUTO_APPROVED" : "CREATED";
        auditService.log(req.getId(), action, dto.getRequestedBy(),
                dto.getHostName() + " - " + dto.getCommand() + " [risk=" + risk + "]");

        return toDto(req);
    }

    public DeferredResult<ResponseEntity<ApprovalResponseDto>> waitForResolution(Long id) {
        ApprovalRequest current = repository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Approval not found: " + id));

        DeferredResult<ResponseEntity<ApprovalResponseDto>> deferred =
                new DeferredResult<>(WAIT_TIMEOUT_MS);

        if (current.getStatus() != ApprovalStatus.PENDING) {
            deferred.setResult(ResponseEntity.ok(toDto(current)));
            return deferred;
        }

        waiters.put(id, deferred);

        deferred.onTimeout(() -> {
            waiters.remove(id);
            ApprovalRequest latest = repository.findById(id).orElseThrow();
            deferred.setResult(ResponseEntity.ok(toDto(latest)));
        });

        deferred.onCompletion(() -> waiters.remove(id));

        return deferred;
    }

    public ApprovalResponseDto resolve(Long id, ApprovalStatus newStatus, String resolvedBy) {
        ApprovalRequest req = repository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Approval not found: " + id));

        req.setStatus(newStatus);
        req.setResolvedAt(Instant.now());
        req.setResolvedBy(resolvedBy);
        repository.save(req);

        auditService.log(req.getId(), newStatus.name(), resolvedBy,
                req.getHostName() + " - " + req.getCommand());

        ApprovalResponseDto dto = toDto(req);

        DeferredResult<ResponseEntity<ApprovalResponseDto>> waiting = waiters.remove(id);
        if (waiting != null) {
            waiting.setResult(ResponseEntity.ok(dto));
        }

        return dto;
    }

    private ApprovalResponseDto toDto(ApprovalRequest req) {
        return ApprovalResponseDto.builder()
                .id(req.getId())
                .requestedBy(req.getRequestedBy())
                .hostName(req.getHostName())
                .command(req.getCommand())
                .riskLevel(req.getRiskLevel())
                .status(req.getStatus())
                .createdAt(req.getCreatedAt())
                .resolvedAt(req.getResolvedAt())
                .resolvedBy(req.getResolvedBy())
                .build();
    }
}