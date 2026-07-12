package com.gozcu.backend.netmon;

import com.gozcu.backend.netmon.dto.NetworkAlertDto;
import com.gozcu.backend.netmon.dto.WhitelistEntryDto;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/netmon")
@RequiredArgsConstructor
public class NetmonController {

    private final NetmonService netmonService;

    @GetMapping("/whitelist")
    public ResponseEntity<List<WhitelistEntryDto>> getWhitelist() {
        return ResponseEntity.ok(netmonService.getWhitelist());
    }

    @PostMapping("/alerts")
    public ResponseEntity<NetworkAlert> createAlert(@RequestBody @Valid NetworkAlertDto dto) {
        return ResponseEntity.ok(netmonService.saveAlert(dto));
    }

    @GetMapping("/alerts")
    public ResponseEntity<List<NetworkAlert>> getAlerts() {
        return ResponseEntity.ok(netmonService.getAlerts());
    }
}
