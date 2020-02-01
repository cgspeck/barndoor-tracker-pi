package wireless

const hostAPDFn = "/etc/hostapd/hostapd.conf"

type hostAPDConfigVars struct {
	Channel   int
	Interface string
	Key       string
	SSID      string
}

const hostAPDConfigTemplate = `interface={{.Interface}}
# hostapd event logger configuration
#
# Two output method: syslog and stdout (only usable if not forking to
# background).
#
# Module bitfield (ORed bitfield of modules that will be logged; -1 = all
# modules):
# bit 0 (1) = IEEE 802.11
# bit 1 (2) = IEEE 802.1X
# bit 2 (4) = RADIUS
# bit 3 (8) = WPA
# bit 4 (16) = driver interface
# bit 5 (32) = IAPP
# bit 6 (64) = MLME
#
# Levels (minimum value for logged events):
#  0 = verbose debugging
#  1 = debugging
#  2 = informational messages
#  3 = notification
#  4 = warning
#
logger_syslog=-1
logger_syslog_level=2
logger_stdout=-1
logger_stdout_level=2

# Interface for separate control program. If this is specified, hostapd
# will create this directory and a UNIX domain socket for listening to requests
# from external programs (CLI/GUI, etc.) for status information and
# configuration. The socket file will be named based on the interface name, so
# multiple hostapd processes/interfaces can be run at the same time if more
# than one interface is used.
# /run/hostapd is the recommended directory for sockets and by default,
# hostapd_cli will use it when trying to connect with hostapd.
ctrl_interface=/run/hostapd

# Access control for the control interface can be configured by setting the
# directory to allow only members of a group to use sockets. This way, it is
# possible to run hostapd as root (since it needs to change network
# configuration and open raw sockets) and still allow GUI/CLI components to be
# run as non-root users. However, since the control interface can be used to
# change the network configuration, this access needs to be protected in many
# cases. By default, hostapd is configured to use gid 0 (root). If you
# want to allow non-root users to use the contron interface, add a new group
# and change this value to match with that group. Add users that should have
# control interface access to this group.
#
# This variable can be a group name or gid.
#ctrl_interface_group=wheel
ctrl_interface_group=0


##### IEEE 802.11 related configuration #######################################

# SSID to be used in IEEE 802.11 management frames
ssid={{.SSID}}

# Country code (ISO/IEC 3166-1). Used to set regulatory domain.
# Set as needed to indicate country in which device is operating.
# This can limit available channels and transmit power.
# These two octets are used as the first two octets of the Country String
# (dot11CountryString)
country_code=AU

# Operation mode (a = IEEE 802.11a (5 GHz), b = IEEE 802.11b (2.4 GHz),
# g = IEEE 802.11g (2.4 GHz), ad = IEEE 802.11ad (60 GHz); a/g options are used
# with IEEE 802.11n (HT), too, to specify band). For IEEE 802.11ac (VHT), this
# needs to be set to hw_mode=a. When using ACS (see channel parameter), a
# special value "any" can be used to indicate that any support band can be used.
# This special case is currently supported only with drivers with which
# offloaded ACS is used.
# Default: IEEE 802.11b
hw_mode=g

# Channel number (IEEE 802.11)
# (default: 0, i.e., not set)
# Please note that some drivers do not use this value from hostapd and the
# channel will need to be configured separately with iwconfig.
#
# If CONFIG_ACS build option is enabled, the channel can be selected
# automatically at run time by setting channel=acs_survey or channel=0, both of
# which will enable the ACS survey based algorithm.
channel={{.Channel}}

# ACS tuning - Automatic Channel Selection
# See: http://wireless.kernel.org/en/users/Documentation/acs
#
# You can customize the ACS survey algorithm with following variables:
#
# acs_num_scans requirement is 1..100 - number of scans to be performed that
# are used to trigger survey data gathering of an underlying device driver.
# Scans are passive and typically take a little over 100ms (depending on the
# driver) on each available channel for given hw_mode. Increasing this value
# means sacrificing startup time and gathering more data wrt channel
# interference that may help choosing a better channel. This can also help fine
# tune the ACS scan time in case a driver has different scan dwell times.
#
# acs_chan_bias is a space-separated list of <channel>:<bias> pairs. It can be
# used to increase (or decrease) the likelihood of a specific channel to be
# selected by the ACS algorithm. The total interference factor for each channel
# gets multiplied by the specified bias value before finding the channel with
# the lowest value. In other words, values between 0.0 and 1.0 can be used to
# make a channel more likely to be picked while values larger than 1.0 make the
# specified channel less likely to be picked. This can be used, e.g., to prefer
# the commonly used 2.4 GHz band channels 1, 6, and 11 (which is the default
# behavior on 2.4 GHz band if no acs_chan_bias parameter is specified).
#
# Defaults:
#acs_num_scans=5
#acs_chan_bias=1:0.8 6:0.8 11:0.8

# Beacon interval in kus (1.024 ms) (default: 100; range 15..65535)
beacon_int=100

# DTIM (delivery traffic information message) period (range 1..255):
# number of beacons between DTIMs (1 = every beacon includes DTIM element)
# (default: 2)
dtim_period=2

# Maximum number of stations allowed in station table. New stations will be
# rejected after the station table is full. IEEE 802.11 has a limit of 2007
# different association IDs, so this number should not be larger than that.
# (default: 2007)
max_num_sta=255

# RTS/CTS threshold; -1 = disabled (default); range -1..65535
# If this field is not included in hostapd.conf, hostapd will not control
# RTS threshold and 'iwconfig wlan# rts <val>' can be used to set it.
rts_threshold=-1

# Fragmentation threshold; -1 = disabled (default); range -1, 256..2346
# If this field is not included in hostapd.conf, hostapd will not control
# fragmentation threshold and 'iwconfig wlan# frag <val>' can be used to set
# it.
fragm_threshold=-1

# Send empty SSID in beacons and ignore probe request frames that do not
# specify full SSID, i.e., require stations to know SSID.
# default: disabled (0)
# 1 = send empty (length=0) SSID in beacon and ignore probe request for
#     broadcast SSID
# 2 = clear SSID (ASCII 0), but keep the original length (this may be required
#     with some clients that do not support empty SSID) and ignore probe
#     requests for broadcast SSID
ignore_broadcast_ssid=0

# Default WMM parameters (IEEE 802.11 draft; 11-03-0504-03-000e):
# for 802.11a or 802.11g networks
# These parameters are sent to WMM clients when they associate.
# The parameters will be used by WMM clients for frames transmitted to the
# access point.
#
# note - txop_limit is in units of 32microseconds
# note - acm is admission control mandatory flag. 0 = admission control not
# required, 1 = mandatory
# note - Here cwMin and cmMax are in exponent form. The actual cw value used
# will be (2^n)-1 where n is the value given here. The allowed range for these
# wmm_ac_??_{cwmin,cwmax} is 0..15 with cwmax >= cwmin.
#
wmm_enabled=1
#
# WMM-PS Unscheduled Automatic Power Save Delivery [U-APSD]
# Enable this flag if U-APSD supported outside hostapd (eg., Firmware/driver)
#uapsd_advertisement_enabled=1
#
# Low priority / AC_BK = background
wmm_ac_bk_cwmin=4
wmm_ac_bk_cwmax=10
wmm_ac_bk_aifs=7
wmm_ac_bk_txop_limit=0
wmm_ac_bk_acm=0
# Note: for IEEE 802.11b mode: cWmin=5 cWmax=10
#
# Normal priority / AC_BE = best effort
wmm_ac_be_aifs=3
wmm_ac_be_cwmin=4
wmm_ac_be_cwmax=10
wmm_ac_be_txop_limit=0
wmm_ac_be_acm=0
# Note: for IEEE 802.11b mode: cWmin=5 cWmax=7
#
# High priority / AC_VI = video
wmm_ac_vi_aifs=2
wmm_ac_vi_cwmin=3
wmm_ac_vi_cwmax=4
wmm_ac_vi_txop_limit=94
wmm_ac_vi_acm=0
# Note: for IEEE 802.11b mode: cWmin=4 cWmax=5 txop_limit=188
#
# Highest priority / AC_VO = voice
wmm_ac_vo_aifs=2
wmm_ac_vo_cwmin=2
wmm_ac_vo_cwmax=3
wmm_ac_vo_txop_limit=47
wmm_ac_vo_acm=0
# Note: for IEEE 802.11b mode: cWmin=3 cWmax=4 burst=102

# EAPOL-Key index workaround (set bit7) for WinXP Supplicant (needed only if
# only broadcast keys are used)
eapol_key_index_workaround=0


##### WPA/IEEE 802.11i configuration ##########################################

# IEEE 802.11 specifies two authentication algorithms. hostapd can be
# configured to allow both of these or only one. Open system authentication
# should be used with IEEE 802.1X.
# Bit fields of allowed authentication algorithms:
# bit 0 = Open System Authentication
# bit 1 = Shared Key Authentication (requires WEP)
auth_algs=1

# Enable WPA. Setting this variable configures the AP to require WPA (either
# WPA-PSK or WPA-RADIUS/EAP based on other configuration). For WPA-PSK, either
# wpa_psk or wpa_passphrase must be set and wpa_key_mgmt must include WPA-PSK.
# Instead of wpa_psk / wpa_passphrase, wpa_psk_radius might suffice.
# For WPA-RADIUS/EAP, ieee8021x must be set (but without dynamic WEP keys),
# RADIUS authentication server must be configured, and WPA-EAP must be included
# in wpa_key_mgmt.
# This field is a bit field that can be used to enable WPA (IEEE 802.11i/D3.0)
# and/or WPA2 (full IEEE 802.11i/RSN):
# bit0 = WPA
# bit1 = IEEE 802.11i/RSN (WPA2) (dot11RSNAEnabled)
# Note that WPA3 is also configured with bit1 since it uses RSN just like WPA2.
# In other words, for WPA3, wpa=2 is used the configuration (and
# wpa_key_mgmt=SAE for WPA3-Personal instead of wpa_key_mgmt=WPA-PSK).
wpa=2
wpa_key_mgmt=WPA-PSK
wpa_passphrase={{.Key}}

# Set of accepted key management algorithms (WPA-PSK, WPA-EAP, or both). The
# entries are separated with a space. WPA-PSK-SHA256 and WPA-EAP-SHA256 can be
# added to enable SHA256-based stronger algorithms.
# WPA-PSK = WPA-Personal / WPA2-Personal
# WPA-PSK-SHA256 = WPA2-Personal using SHA256
# WPA-EAP = WPA-Enterprise / WPA2-Enterprise
# WPA-EAP-SHA256 = WPA2-Enterprise using SHA256
# SAE = SAE (WPA3-Personal)
# WPA-EAP-SUITE-B-192 = WPA3-Enterprise with 192-bit security/CNSA suite
# FT-PSK = FT with passphrase/PSK
# FT-EAP = FT with EAP
# FT-EAP-SHA384 = FT with EAP using SHA384
# FT-SAE = FT with SAE
# FILS-SHA256 = Fast Initial Link Setup with SHA256
# FILS-SHA384 = Fast Initial Link Setup with SHA384
# FT-FILS-SHA256 = FT and Fast Initial Link Setup with SHA256
# FT-FILS-SHA384 = FT and Fast Initial Link Setup with SHA384
# OWE = Opportunistic Wireless Encryption (a.k.a. Enhanced Open)
# DPP = Device Provisioning Protocol
# OSEN = Hotspot 2.0 online signup with encryption
# (dot11RSNAConfigAuthenticationSuitesTable)
wpa_key_mgmt=WPA-PSK

##### RADIUS client configuration #############################################
# for IEEE 802.1X with external Authentication Server, IEEE 802.11
# authentication with external ACL for MAC addresses, and accounting

# The own IP address of the access point (used as NAS-IP-Address)
own_ip_addr=127.0.0.1

##### Wi-Fi Protected Setup (WPS) #############################################

# WPS state
# 0 = WPS disabled (default)
# 1 = WPS enabled, not configured
# 2 = WPS enabled, configured
#wps_state=2

# Whether to manage this interface independently from other WPS interfaces
# By default, a single hostapd process applies WPS operations to all configured
# interfaces. This parameter can be used to disable that behavior for a subset
# of interfaces. If this is set to non-zero for an interface, WPS commands
# issued on that interface do not apply to other interfaces and WPS operations
# performed on other interfaces do not affect this interface.
#wps_independent=0

# AP can be configured into a locked state where new WPS Registrar are not
# accepted, but previously authorized Registrars (including the internal one)
# can continue to add new Enrollees.
#ap_setup_locked=1

# Universally Unique IDentifier (UUID; see RFC 4122) of the device
# This value is used as the UUID for the internal WPS Registrar. If the AP
# is also using UPnP, this value should be set to the device's UPnP UUID.
# If not configured, UUID will be generated based on the local MAC address.
#uuid=12345678-9abc-def0-1234-56789abcdef0

# Note: If wpa_psk_file is set, WPS is used to generate random, per-device PSKs
# that will be appended to the wpa_psk_file. If wpa_psk_file is not set, the
# default PSK (wpa_psk/wpa_passphrase) will be delivered to Enrollees. Use of
# per-device PSKs is recommended as the more secure option (i.e., make sure to
# set wpa_psk_file when using WPS with WPA-PSK).

# When an Enrollee requests access to the network with PIN method, the Enrollee
# PIN will need to be entered for the Registrar. PIN request notifications are
# sent to hostapd ctrl_iface monitor. In addition, they can be written to a
# text file that could be used, e.g., to populate the AP administration UI with
# pending PIN requests. If the following variable is set, the PIN requests will
# be written to the configured file.
#wps_pin_requests=/run/hostapd/wps_pin_requests

# Device Name
# User-friendly description of device; up to 32 octets encoded in UTF-8
#device_name=Wireless AP

# Manufacturer
# The manufacturer of the device (up to 64 ASCII characters)
#manufacturer=Company

# Model Name
# Model of the device (up to 32 ASCII characters)
#model_name=WAP

# Model Number
# Additional device description (up to 32 ASCII characters)
#model_number=123

# Serial Number
# Serial number of the device (up to 32 characters)
#serial_number=12345

# Primary Device Type
# Used format: <categ>-<OUI>-<subcateg>
# categ = Category as an integer value
# OUI = OUI and type octet as a 4-octet hex-encoded value; 0050F204 for
#       default WPS OUI
# subcateg = OUI-specific Sub Category as an integer value
# Examples:
#   1-0050F204-1 (Computer / PC)
#   1-0050F204-2 (Computer / Server)
#   5-0050F204-1 (Storage / NAS)
#   6-0050F204-1 (Network Infrastructure / AP)
#device_type=6-0050F204-1

# OS Version
# 4-octet operating system version number (hex string)
#os_version=01020300

# Config Methods
# List of the supported configuration methods
# Available methods: usba ethernet label display ext_nfc_token int_nfc_token
#	nfc_interface push_button keypad virtual_display physical_display
#	virtual_push_button physical_push_button
#config_methods=label virtual_display virtual_push_button keypad

# WPS capability discovery workaround for PBC with Windows 7
# Windows 7 uses incorrect way of figuring out AP's WPS capabilities by acting
# as a Registrar and using M1 from the AP. The config methods attribute in that
# message is supposed to indicate only the configuration method supported by
# the AP in Enrollee role, i.e., to add an external Registrar. For that case,
# PBC shall not be used and as such, the PushButton config method is removed
# from M1 by default. If pbc_in_m1=1 is included in the configuration file,
# the PushButton config method is left in M1 (if included in config_methods
# parameter) to allow Windows 7 to use PBC instead of PIN (e.g., from a label
# in the AP).
#pbc_in_m1=1

# Static access point PIN for initial configuration and adding Registrars
# If not set, hostapd will not allow external WPS Registrars to control the
# access point. The AP PIN can also be set at runtime with hostapd_cli
# wps_ap_pin command. Use of temporary (enabled by user action) and random
# AP PIN is much more secure than configuring a static AP PIN here. As such,
# use of the ap_pin parameter is not recommended if the AP device has means for
# displaying a random PIN.
#ap_pin=12345670

# Skip building of automatic WPS credential
# This can be used to allow the automatically generated Credential attribute to
# be replaced with pre-configured Credential(s).
#skip_cred_build=1

# Additional Credential attribute(s)
# This option can be used to add pre-configured Credential attributes into M8
# message when acting as a Registrar. If skip_cred_build=1, this data will also
# be able to override the Credential attribute that would have otherwise been
# automatically generated based on network configuration. This configuration
# option points to an external file that much contain the WPS Credential
# attribute(s) as binary data.
#extra_cred=/var/lib/hostapd/hostapd.cred

# Credential processing
#   0 = process received credentials internally (default)
#   1 = do not process received credentials; just pass them over ctrl_iface to
#	external program(s)
#   2 = process received credentials internally and pass them over ctrl_iface
#	to external program(s)
# Note: With wps_cred_processing=1, skip_cred_build should be set to 1 and
# extra_cred be used to provide the Credential data for Enrollees.
#
# wps_cred_processing=1 will disabled automatic updates of hostapd.conf file
# both for Credential processing and for marking AP Setup Locked based on
# validation failures of AP PIN. An external program is responsible on updating
# the configuration appropriately in this case.
#wps_cred_processing=0

# Whether to enable SAE (WPA3-Personal transition mode) automatically for
# WPA2-PSK credentials received using WPS.
# 0 = only add the explicitly listed WPA2-PSK configuration (default)
# 1 = add both the WPA2-PSK and SAE configuration and enable PMF so that the
#     AP gets configured in WPA3-Personal transition mode (supports both
#     WPA2-Personal (PSK) and WPA3-Personal (SAE) clients).
#wps_cred_add_sae=0

# AP Settings Attributes for M7
# By default, hostapd generates the AP Settings Attributes for M7 based on the
# current configuration. It is possible to override this by providing a file
# with pre-configured attributes. This is similar to extra_cred file format,
# but the AP Settings attributes are not encapsulated in a Credential
# attribute.
#ap_settings=/var/lib/hostapd/hostapd.ap_settings

# Multi-AP backhaul BSS config
# Used in WPS when multi_ap=2 or 3. Defines "backhaul BSS" credentials.
# These are passed in WPS M8 instead of the normal (fronthaul) credentials
# if the Enrollee has the Multi-AP subelement set. Backhaul SSID is formatted
# like ssid2. The key is set like wpa_psk or wpa_passphrase.
#multi_ap_backhaul_ssid="backhaul"
#multi_ap_backhaul_wpa_psk=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
#multi_ap_backhaul_wpa_passphrase=secret passphrase

# WPS UPnP interface
# If set, support for external Registrars is enabled.
#upnp_iface=br0

# Friendly Name (required for UPnP)
# Short description for end use. Should be less than 64 characters.
#friendly_name=WPS Access Point

# Manufacturer URL (optional for UPnP)
#manufacturer_url=http://www.example.com/

# Model Description (recommended for UPnP)
# Long description for end user. Should be less than 128 characters.
#model_description=Wireless Access Point

# Model URL (optional for UPnP)
#model_url=http://www.example.com/model/

# Universal Product Code (optional for UPnP)
# 12-digit, all-numeric code that identifies the consumer package.
#upc=123456789012

# WPS RF Bands (a = 5G, b = 2.4G, g = 2.4G, ag = dual band, ad = 60 GHz)
# This value should be set according to RF band(s) supported by the AP if
# hw_mode is not set. For dual band dual concurrent devices, this needs to be
# set to ag to allow both RF bands to be advertized.
#wps_rf_bands=ag

# NFC password token for WPS
# These parameters can be used to configure a fixed NFC password token for the
# AP. This can be generated, e.g., with nfc_pw_token from wpa_supplicant. When
# these parameters are used, the AP is assumed to be deployed with a NFC tag
# that includes the matching NFC password token (e.g., written based on the
# NDEF record from nfc_pw_token).
#
#wps_nfc_dev_pw_id: Device Password ID (16..65535)
#wps_nfc_dh_pubkey: Hexdump of DH Public Key
#wps_nfc_dh_privkey: Hexdump of DH Private Key
#wps_nfc_dev_pw: Hexdump of Device Password
`
