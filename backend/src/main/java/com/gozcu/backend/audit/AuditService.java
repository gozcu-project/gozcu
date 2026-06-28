package com.gozcu.backend.audit;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import java.time.Instant;

@Service
@RequiredArgsConstructor
public class AuditService {

    private final AuditLogRepository repository;

    public void log(Long approvalRequestId, String action, String actor, String details) {
        AuditLog entry = AuditLog.builder()
                .approvalRequestId(approvalRequestId)
                .action(action)
                .actor(actor)
                .timestamp(Instant.now())
                .details(details)
                .build();
        repository.save(entry);
    }
}