# Gözcü & Seren Linux — Yol Haritası

Bu belge, Gözcü'nün açık kaynak Linux Security Platform olarak gelişim sürecini ve
uzun vadeli hedef olan Seren Linux dağıtımını kapsar. Süreler kesin taahhüt değil,
tek geliştirici + part-time çalışma varsayımıyla yapılmış gerçekçi tahminlerdir.
Topluluk katkısı bu süreleri önemli ölçüde kısaltabilir.

---

## Faz A — Spring Boot Backend ✅ Tamamlandı
- Approval, Audit, Fleet, Policy modülleri
- Hibrit DeferredResult tabanlı onay bekleme (30+30sn)
- Risk skoru bazlı otomatik onay (LOW → anında geçiş)
- PostgreSQL 18 entegrasyonu

## Faz B — mTLS ✅ Tamamlandı
- Kendi CA'sı (OpenSSL)
- Server/client sertifikaları (SAN dahil)
- Spring Boot SSL client-auth (client-auth: need)
- Go tarafında crypto/tls entegrasyonu

## Faz C — gozcu-gate ✅ Tamamlandı
- sudoers tabanlı komut yakalama (gate-only erişim)
- Go binary (statik, CGO_ENABLED=0)
- mTLS client ile backend'e POST
- Onay sonrası exec() ile gerçek komuta geçiş
- .bashrc fonksiyon enjeksiyonu (UX katmanı)
- Uçtan uca test: sudo → gate → mTLS → backend → onay → exec

---

## Faz D — gozcu-netmon (eBPF Network Monitor)
- libbpf-go + CO-RE/BTF tabanlı (kernel-bağımsız)
- connect/socket syscall hook'ları
- Whitelist tabanlı network policy (backend'den çekilen)
- Whitelist dışı bağlantılarda mTLS üzerinden alarm
- Ubuntu VPS'te CO-RE doğrulaması (Fedora'da geliştirme)
- gozcu-netmon systemd servisi olarak kurulum

## Faz E — Angular Arayüzü
- WebSocket/STOMP ile gerçek zamanlı approval kartı
- YES/NO onay butonu
- Whitelist CRUD ekranı
- Alert dashboard
- Fleet listesi (host yönetimi)
- Policy CRUD ekranı (risk seviyesi yönetimi)

## Faz F — Entegrasyon ve Kapanış
- gozcu-gate + gozcu-netmon + backend + Angular uçtan uca test
- Docker Compose ile tam deployment
- README güncellemesi (demo kaydı, kurulum kılavuzu)
- v0.1.0 GitHub release

---

## Faz 1 — Configuration Security
- Kritik sistem dosyalarının (sshd_config, sudoers, crontab vb.) değişiklik tespiti
- Değişiklik öncesi onay mekanizmasıyla entegrasyon
- Baseline snapshot ve drift tespiti

## Faz 2 — Audit Derinleştirme
- Audit log'a TIMEOUT kaydı (✅ eklendi)
- Audit log export (CSV/JSON)
- Audit log arama ve filtreleme
- Uzun süreli audit saklama politikası

## Faz 3 — Fleet Management
- Host onboarding scripti (otomatik sudoers + gate kurulumu)
- Host health-check: sudoers drift tespiti (kullanıcı %sudo grubuna geri eklendiyse uyarı)
- Host grupları ve grup bazlı policy yönetimi
- Online/offline host durumu takibi

## Faz 4 — SSL/TLS Sertifika Yönetimi
- Sistemdeki SSL/TLS, JKS, PKCS#12 sertifikalarının merkezi takibi
- Sertifika son kullanma tarihi uyarıları
- Imzala.app entegrasyonu (sertifika yenileme akışı)

## Faz 5 — Infrastructure Monitoring (Temel)
- CPU, RAM, disk kullanım metrikleri (agent tabanlı)
- Eşik aşımında alarm üretimi
- Backend'e metrik gönderimi, Angular'da basit dashboard

## Faz 6 — JVM Diagnostics
- Spring Boot/Java tabanlı uygulamalar için JVM metrik izleme
- Heap, GC, thread durumu takibi
- Actuator entegrasyonu

## Faz 7 — Log Analytics
- Sistem log'larının (syslog, journald) merkezi toplanması
- Pattern bazlı alarm kuralları
- Log arama ve filtreleme arayüzü

## Faz 8 — Database Diagnostics
- PostgreSQL, MySQL bağlantı havuzu ve sorgu izleme
- Yavaş sorgu tespiti ve alarm
- DB health-check entegrasyonu

## Faz 9 — Deploy Guard
- CI/CD pipeline entegrasyonu (GitHub Actions, GitLab CI)
- Deploy öncesi onay mekanizması
- Yetkisiz deploy engelleme

## Faz 10 — SSH Security
- SSH login izleme (başarılı/başarısız)
- Brute-force tespiti
- Yetkisiz SSH key tespiti

## Faz 11 — MFA (Multi-Factor Authentication)
- TOTP tabanlı ikinci faktör (Google Authenticator uyumlu)
- MEDIUM/HIGH risk komutlar için MFA zorunluluğu
- MFA bypass audit kaydı

## Faz 12 — Incident Timeline
- Bir olayın başından sonuna tüm adımlarının kronolojik görselleştirmesi
- "Bu saldırı nasıl başladı, hangi adımlardan geçti?" sorusuna cevap
- Approval, audit, alert, network event'larının tek zaman çizelgesinde birleştirilmesi

## Faz 13 — Root Cause Analysis
- Otomatik kök neden analizi (hangi komut, hangi kullanıcı, hangi host, ne zaman)
- Bağlam zinciri: önceki olaylarla ilişkilendirme
- Öneri motoru: "bir daha olmaması için şunu yapın"

## Faz 14 — Slack/Teams Entegrasyonu
- Onay isteklerinin Slack/Teams kanalına gönderilmesi
- Kanal üzerinden YES/NO onayı
- Alert bildirimlerinin Slack/Teams'e yönlendirilmesi

## Faz 15 — Smart/Context-Aware Approval
- Onay kararında bağlamsal faktörlerin dikkate alınması
  (mesai saati mi? production mu? ilk kez mi çalıştırılıyor?)
- Dinamik risk skoru hesaplama
- "Bu koşullarda bu komut için onay gereksiz" kuralları

## Faz 16 — Security Analytics
- Kullanıcı/host bazlı güvenlik skoru
- Trend analizi: risk seviyesi değişimleri
- Anomali tespitine giriş: alışılmadık saatlerde komut çalıştırma

## Faz 17 — Sudoers Self-Integrity Enforcement
- systemd timer ile periyodik sudoers drift kontrolü
- Kullanıcı %sudo grubuna geri eklendiyse otomatik çıkarma
- Enforcement log'u audit'e yazma
- "Sadece root gate'i kırabilir" garantisinin otomatik korunması

---

## Faz 18 — MITRE ATT&CK Entegrasyonu
- Policy tablosuna ATT&CK TTP referans alanı eklenmesi
  (her komut pattern'ine Tactic/Technique/Sub-technique eşleştirmesi)
- Audit log'da "bu işlem hangi ATT&CK taktiğine karşılık geliyor" bilgisi
- Angular'da ATT&CK bazlı policy yönetimi ekranı
- Hedef: threat intelligence'a dayalı policy yönetimi

Örnek eşleştirmeler:
| Komut | ATT&CK TTP |
|---|---|
| chmod 777 * | T1222 — File Permission Modification |
| crontab * | T1053 — Scheduled Task/Job |
| iptables -F | T1562 — Impair Defenses |
| useradd/usermod * | T1136 — Create Account |
| visudo | T1548 — Abuse Elevation Control Mechanism |
| systemctl disable * | T1562 — Impair Defenses |

## Faz 19 — Tehdit Korelasyon Motoru
- Tek tek sudo olaylarını değil, zaman penceresi içindeki olay zincirlerini analiz etme
- "Aynı kullanıcı 5 dakika içinde 3 farklı MEDIUM risk komut denedi" örüntüsü tespiti
- Saldırı zinciri görselleştirmesi (Faz 12 Incident Timeline ile entegrasyon)
- SIEM'in yaptığı korelasyonu privileged access katmanında yerel olarak yapabilmek

## Faz 20 — Davranışsal Anomali Tespiti (Behavioral Baseline)
- Her kullanıcı/host için normal davranış profili oluşturma
- Baseline'dan sapmaların otomatik risk skoru artışına yansıtılması
  ("bu kullanıcı bu komutu hiç çalıştırmamış" → otomatik HIGH)
- Statik policy tablosunun ötesine geçip dinamik/adaptif risk değerlendirmesi

## Faz 21 — Entegrasyon Ekosistemi
- SIEM entegrasyonu: Splunk, Elastic, IBM QRadar
  (audit loglarını doğrudan SIEM'e gönderme)
- EDR entegrasyonu: CrowdStrike, SentinelOne, Microsoft Defender
  (EDR'dan gelen tehdit sinyallerini Gözcü'nün risk skoruna yansıtma)
- Webhook/API gateway: herhangi bir güvenlik aracının Gözcü event'larını tüketebilmesi
- Hedef: izole bir ada olmaktan çıkıp mevcut MDR ekosistemiyle konuşan bir bileşen

## Faz 22 — Incident Response Desteği
- Aktif saldırı anında yarı-otomatik müdahale
  ("şu host'u izole et, şu kullanıcıyı geçici kilitle")
- Playbook desteği: belirli saldırı örüntüsü tespit edilince önceden tanımlı adımları otomatik çalıştırma
- Break-glass mekanizması: acil durumda yetkili kişinin kısıtlamaları geçici kaldırabilmesi
  (ve bunun audit'e eksiksiz düşmesi)
- Hedef: "tespit et ve raporla"dan "tespit et ve müdahale et"e geçiş

## Faz 23 — Raporlama ve MDR Dashboard
- MTTD (Mean Time to Detect) ve MTTR (Mean Time to Respond) metrikleri
- Kullanıcı/host bazlı risk trend analizi
- Yönetici raporları (PDF/export): üst yönetime veya müşterilere sunulabilecek format
- Hedef: MDR ekip liderinin haftalık/aylık raporlamasını Gözcü üzerinden yapabilmesi

## Faz 24 — Filesystem Activity Monitor (Ransomware Erken Tespiti)
- eBPF ile dosya syscall izleme (openat, write, rename, unlink)
- Toplu şifreleme paterni tespiti (kısa sürede çok sayıda dosya değişimi)
- Canary/honeypot dosya tuzakları
- Önemli not: Bu bileşen tespit ve erken uyarı sağlar, şifrelemeyi engellemez.
  Gerçek koruma immutable/offsite backup (S3 Object Lock gibi) ile sağlanır.

## Faz 25 — Supply Chain Security
- Git repository izleme: yetkisiz commit/push tespiti
- Paket bütünlük kontrolü (npm, pip, maven bağımlılıkları)
- CI/CD pipeline anomali tespiti
- Bağımlılık zinciri görselleştirmesi

## Faz 26 — Enterprise (K8s/Docker/LDAP/SSO)
- Kubernetes entegrasyonu: pod seviyesinde komut onayı
- Docker entegrasyonu: container exec onayı
- LDAP/Active Directory entegrasyonu
- SSO (SAML/OIDC) desteği

---

## Seren Linux

### Vizyon
Gözcü'nün tüm güvenlik bileşenlerinin çekirdekte yerleşik geldiği,
kullanıcı ve grup yönetiminin dağıtımın kendisi tarafından dayatıldığı,
sudoers yanlış konfigürasyonunun yapısal olarak mümkün olmadığı
bir Linux dağıtımı.

Güvenlik, sonradan eklenen bir katman değil — sistemin ve çekirdeğin bizzat kendisi.

### Paket Desteği
- .deb (Debian/Ubuntu tabanlı) — geniş topluluk, yaygın kurumsal kullanım
- .rpm (RHEL/Rocky/AlmaLinux tabanlı) — kurumsal/kamu sektörü (Türkiye'de yaygın)

### Entegre Ürünler
- **Gözcü** — privileged access management, command approval, audit
- **Imzala.app** — eIDAS uyumlu dijital imza, inkâr edilemez kayıt

### Uyumluluk Hedefleri
- ISO 27001/27002 (privileged access management, erişim kaydı)
- KVKK (erişim denetimi yükümlülükleri)
- eIDAS (inkâr edilemezlik — Imzala.app üzerinden)

### Geliştirme Tahmini (tek geliştirici, part-time)
| Milestone | Tahmini Süre |
|---|---|
| Gözcü Faz D-F tamamlanması | 3-4 ay |
| Gözcü Faz 1-17 tamamlanması | 6-9 ay |
| Gözcü Faz 18-26 tamamlanması | 9-12 ay |
| Seren Linux v0.1 (beta dağıtım) | +3-6 ay |
| **Toplam: Seren Linux v0.1** | **~18-24 ay** |

Topluluk katkısı bu süreleri önemli ölçüde kısaltabilir.
Her fazın çıktısı bağımsız olarak yayınlanır — dağıtım son halka, ama
önceki her halka kendi başına değerlidir.

### Tehdit Modeli Notu
Root'u kısıtlamak kobra etkisi yaratır. Seren Linux'un hedefi root'u
ortadan kaldırmak değil, root yetkisinin kullanımını görünür, hesap verebilir
ve onay mekanizmasına bağlı kılmaktır. Gerçek root erişimi (fiziksel, kernel exploit)
her zaman kapsam dışıdır — bu bir eksiklik değil, bilinçli bir tasarım sınırıdır.

---

*Bu belge yaşayan bir dokümandır. Her fazın tamamlanmasıyla güncellenir.*
*Son güncelleme: Temmuz 2026*
