package com.gozcu.backend.approval;

import com.gozcu.backend.approval.ApprovalRequest;

import org.springframework.data.jpa.repository.JpaRepository;

public interface ApprovalRequestRepository extends JpaRepository<ApprovalRequest, Long> {
}