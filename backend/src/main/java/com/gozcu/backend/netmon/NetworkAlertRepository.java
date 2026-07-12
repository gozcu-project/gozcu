package com.gozcu.backend.netmon;

import org.springframework.data.jpa.repository.JpaRepository;
import java.util.List;

public interface NetworkAlertRepository extends JpaRepository<NetworkAlert, Long> {
    List<NetworkAlert> findAllByOrderByDetectedAtDesc();
}
