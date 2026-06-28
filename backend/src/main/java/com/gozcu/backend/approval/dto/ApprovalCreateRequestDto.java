package com.gozcu.backend.approval.dto;

import jakarta.validation.constraints.NotBlank;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class ApprovalCreateRequestDto {

    @NotBlank
    private String requestedBy;

    @NotBlank
    private String hostName;

    @NotBlank
    private String command;
}