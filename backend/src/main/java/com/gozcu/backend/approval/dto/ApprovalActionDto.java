package com.gozcu.backend.approval.dto;

import jakarta.validation.constraints.NotBlank;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class ApprovalActionDto {

    @NotBlank
    private String resolvedBy;
}