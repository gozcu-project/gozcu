package com.gozcu.backend.fleet.dto;

import jakarta.validation.constraints.NotBlank;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class HostRegisterDto {

    @NotBlank
    private String hostname;

    @NotBlank
    private String ipAddress;
}