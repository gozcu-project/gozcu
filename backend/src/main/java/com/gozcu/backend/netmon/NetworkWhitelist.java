package com.gozcu.backend.netmon;

import jakarta.persistence.*;
import lombok.*;

@Entity
@Table(name = "network_whitelists")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class NetworkWhitelist {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false)
    private String cidr;

    @Column
    private String ports;

    @Column
    private String description;
}
