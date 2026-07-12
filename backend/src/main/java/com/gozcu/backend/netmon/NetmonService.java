package com.gozcu.backend.netmon;

import com.gozcu.backend.netmon.dto.NetworkAlertDto;
import com.gozcu.backend.netmon.dto.WhitelistEntryDto;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.time.Instant;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class NetmonService {

    private final NetworkWhitelistRepository whitelistRepository;
    private final NetworkAlertRepository alertRepository;

    public List<WhitelistEntryDto> getWhitelist() {
        return whitelistRepository.findAll().stream()
                .map(this::toWhitelistDto)
                .collect(Collectors.toList());
    }

    public NetworkAlert saveAlert(NetworkAlertDto dto) {
        NetworkAlert alert = NetworkAlert.builder()
                .pid(dto.getPid())
                .uid(dto.getUid())
                .comm(dto.getComm())
                .destIp(dto.getDestIp())
                .dstPort(dto.getDstPort())
                .proto(dto.getProto())
                .detectedAt(Instant.now())
                .build();
        return alertRepository.save(alert);
    }

    public List<NetworkAlert> getAlerts() {
        return alertRepository.findAllByOrderByDetectedAtDesc();
    }

    private WhitelistEntryDto toWhitelistDto(NetworkWhitelist entry) {
        List<Integer> ports = (entry.getPorts() == null || entry.getPorts().isBlank())
                ? Collections.emptyList()
                : Arrays.stream(entry.getPorts().split(","))
                        .map(String::trim)
                        .map(Integer::parseInt)
                        .collect(Collectors.toList());

        return WhitelistEntryDto.builder()
                .cidr(entry.getCidr())
                .ports(ports)
                .build();
    }
}
