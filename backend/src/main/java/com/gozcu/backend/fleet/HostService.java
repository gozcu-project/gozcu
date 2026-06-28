package com.gozcu.backend.fleet;


import com.gozcu.backend.fleet.dto.HostRegisterDto;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import java.time.Instant;

@Service
@RequiredArgsConstructor
public class HostService {

    private final HostRepository repository;

    public Host register(HostRegisterDto dto) {
        Host host = Host.builder()
                .hostname(dto.getHostname())
                .ipAddress(dto.getIpAddress())
                .registeredAt(Instant.now())
                .build();
        return repository.save(host);
    }
}