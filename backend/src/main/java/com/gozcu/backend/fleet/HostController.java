package com.gozcu.backend.fleet;


import com.gozcu.backend.fleet.HostService;
import com.gozcu.backend.fleet.dto.HostRegisterDto;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import java.util.List;

@RestController
@RequestMapping("/api/hosts")
@RequiredArgsConstructor
public class HostController {

    private final HostService hostService;
    private final HostRepository hostRepository;

    @PostMapping
    public ResponseEntity<Host> register(@RequestBody @Valid HostRegisterDto dto) {
        return ResponseEntity.ok(hostService.register(dto));
    }

    @GetMapping
    public List<Host> list() {
        return hostRepository.findAll();
    }
}