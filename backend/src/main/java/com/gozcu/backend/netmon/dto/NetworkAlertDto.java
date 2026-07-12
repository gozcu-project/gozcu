package com.gozcu.backend.netmon.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class NetworkAlertDto {

    @NotNull
    private Long pid;

    @NotNull
    private Long uid;

    @NotBlank
    private String comm;

    @NotBlank
    private String destIp;

    @NotNull
    private Integer dstPort;

    @NotBlank
    private String proto;
}
