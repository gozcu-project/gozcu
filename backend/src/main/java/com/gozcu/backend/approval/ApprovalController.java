package com.gozcu.backend.approval;

import com.gozcu.backend.approval.dto.ApprovalActionDto;
import com.gozcu.backend.approval.dto.ApprovalCreateRequestDto;
import com.gozcu.backend.approval.dto.ApprovalResponseDto;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.context.request.async.DeferredResult;

@RestController
@RequestMapping("/api/approvals")
@RequiredArgsConstructor
public class ApprovalController {

    private final ApprovalService approvalService;

    @PostMapping
    public ResponseEntity<ApprovalResponseDto> create(@RequestBody @Valid ApprovalCreateRequestDto dto) {
        return ResponseEntity.ok(approvalService.create(dto));
    }

    @GetMapping("/{id}/wait")
    public DeferredResult<ResponseEntity<ApprovalResponseDto>> wait(@PathVariable Long id) {
        return approvalService.waitForResolution(id);
    }

    @PutMapping("/{id}/approve")
    public ResponseEntity<ApprovalResponseDto> approve(@PathVariable Long id,
                                                       @RequestBody @Valid ApprovalActionDto dto) {
        return ResponseEntity.ok(approvalService.resolve(id, ApprovalStatus.APPROVED, dto.getResolvedBy()));
    }

    @PutMapping("/{id}/reject")
    public ResponseEntity<ApprovalResponseDto> reject(@PathVariable Long id,
                                                      @RequestBody @Valid ApprovalActionDto dto) {
        return ResponseEntity.ok(approvalService.resolve(id, ApprovalStatus.REJECTED, dto.getResolvedBy()));
    }
}