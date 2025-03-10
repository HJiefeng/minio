// Copyright (c) 2015-2021 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/minio/madmin-go"
	"github.com/minio/minio-go/v7/pkg/set"
	"github.com/minio/minio/internal/auth"
	"github.com/minio/pkg/env"
)

// Error config error type
type Error struct {
	Err string
}

// Errorf - formats according to a format specifier and returns
// the string as a value that satisfies error of type config.Error
func Errorf(format string, a ...interface{}) error {
	return Error{Err: fmt.Sprintf(format, a...)}
}

func (e Error) Error() string {
	return e.Err
}

// Default keys
const (
	Default = madmin.Default
	Enable  = madmin.EnableKey
	Comment = madmin.CommentKey

	// Enable values
	EnableOn  = madmin.EnableOn
	EnableOff = madmin.EnableOff

	RegionKey  = "region"
	NameKey    = "name"
	RegionName = "name"
	AccessKey  = "access_key"
	SecretKey  = "secret_key"
	License    = "license" // Deprecated Dec 2021
	APIKey     = "api_key"
	Proxy      = "proxy"
)

// Top level config constants.
const (
	CredentialsSubSys    = "credentials"
	PolicyOPASubSys      = "policy_opa"
	PolicyPluginSubSys   = "policy_plugin"
	IdentityOpenIDSubSys = "identity_openid"
	IdentityLDAPSubSys   = "identity_ldap"
	IdentityTLSSubSys    = "identity_tls"
	IdentityPluginSubSys = "identity_plugin"
	CacheSubSys          = "cache"
	SiteSubSys           = "site"
	RegionSubSys         = "region"
	EtcdSubSys           = "etcd"
	StorageClassSubSys   = "storage_class"
	APISubSys            = "api"
	CompressionSubSys    = "compression"
	LoggerWebhookSubSys  = "logger_webhook"
	AuditWebhookSubSys   = "audit_webhook"
	AuditKafkaSubSys     = "audit_kafka"
	HealSubSys           = "heal"
	ScannerSubSys        = "scanner"
	CrawlerSubSys        = "crawler"
	SubnetSubSys         = "subnet"
	CallhomeSubSys       = "callhome"

	// Add new constants here if you add new fields to config.
)

// Notification config constants.
const (
	NotifyKafkaSubSys    = "notify_kafka"
	NotifyMQTTSubSys     = "notify_mqtt"
	NotifyMySQLSubSys    = "notify_mysql"
	NotifyNATSSubSys     = "notify_nats"
	NotifyNSQSubSys      = "notify_nsq"
	NotifyESSubSys       = "notify_elasticsearch"
	NotifyAMQPSubSys     = "notify_amqp"
	NotifyPostgresSubSys = "notify_postgres"
	NotifyRedisSubSys    = "notify_redis"
	NotifyWebhookSubSys  = "notify_webhook"

	// Add new constants here if you add new fields to config.
)

// NotifySubSystems - all notification sub-systems
var NotifySubSystems = set.CreateStringSet(
	NotifyKafkaSubSys,
	NotifyMQTTSubSys,
	NotifyMySQLSubSys,
	NotifyNATSSubSys,
	NotifyNSQSubSys,
	NotifyESSubSys,
	NotifyAMQPSubSys,
	NotifyPostgresSubSys,
	NotifyRedisSubSys,
	NotifyWebhookSubSys,
)

// LoggerSubSystems - all sub-systems related to logger
var LoggerSubSystems = set.CreateStringSet(
	LoggerWebhookSubSys,
	AuditWebhookSubSys,
	AuditKafkaSubSys,
)

// SubSystems - all supported sub-systems
var SubSystems = set.CreateStringSet(
	CredentialsSubSys,
	SiteSubSys,
	RegionSubSys,
	EtcdSubSys,
	CacheSubSys,
	APISubSys,
	StorageClassSubSys,
	CompressionSubSys,
	LoggerWebhookSubSys,
	AuditWebhookSubSys,
	AuditKafkaSubSys,
	PolicyOPASubSys,
	PolicyPluginSubSys,
	IdentityLDAPSubSys,
	IdentityOpenIDSubSys,
	IdentityTLSSubSys,
	IdentityPluginSubSys,
	ScannerSubSys,
	HealSubSys,
	NotifyAMQPSubSys,
	NotifyESSubSys,
	NotifyKafkaSubSys,
	NotifyMQTTSubSys,
	NotifyMySQLSubSys,
	NotifyNATSSubSys,
	NotifyNSQSubSys,
	NotifyPostgresSubSys,
	NotifyRedisSubSys,
	NotifyWebhookSubSys,
	SubnetSubSys,
	CallhomeSubSys,
)

// SubSystemsDynamic - all sub-systems that have dynamic config.
var SubSystemsDynamic = set.CreateStringSet(
	APISubSys,
	CompressionSubSys,
	ScannerSubSys,
	HealSubSys,
	SubnetSubSys,
	CallhomeSubSys,
	LoggerWebhookSubSys,
	AuditWebhookSubSys,
	AuditKafkaSubSys,
	StorageClassSubSys,
)

// SubSystemsSingleTargets - subsystems which only support single target.
var SubSystemsSingleTargets = set.CreateStringSet([]string{
	CredentialsSubSys,
	SiteSubSys,
	RegionSubSys,
	EtcdSubSys,
	CacheSubSys,
	APISubSys,
	StorageClassSubSys,
	CompressionSubSys,
	PolicyOPASubSys,
	PolicyPluginSubSys,
	IdentityLDAPSubSys,
	IdentityTLSSubSys,
	IdentityPluginSubSys,
	HealSubSys,
	ScannerSubSys,
}...)

// Constant separators
const (
	SubSystemSeparator = madmin.SubSystemSeparator
	KvSeparator        = madmin.KvSeparator
	KvSpaceSeparator   = madmin.KvSpaceSeparator
	KvComment          = madmin.KvComment
	KvNewline          = madmin.KvNewline
	KvDoubleQuote      = madmin.KvDoubleQuote
	KvSingleQuote      = madmin.KvSingleQuote

	// Env prefix used for all envs in MinIO
	EnvPrefix        = "MINIO_"
	EnvWordDelimiter = `_`
)

// DefaultKVS - default kvs for all sub-systems
var DefaultKVS map[string]KVS

// RegisterDefaultKVS - this function saves input kvsMap
// globally, this should be called only once preferably
// during `init()`.
func RegisterDefaultKVS(kvsMap map[string]KVS) {
	DefaultKVS = map[string]KVS{}
	for subSys, kvs := range kvsMap {
		DefaultKVS[subSys] = kvs
	}
}

// HelpSubSysMap - help for all individual KVS for each sub-systems
// also carries a special empty sub-system which dumps
// help for each sub-system key.
var HelpSubSysMap map[string]HelpKVS

// RegisterHelpSubSys - this function saves
// input help KVS for each sub-system globally,
// this function should be called only once
// preferably in during `init()`.
func RegisterHelpSubSys(helpKVSMap map[string]HelpKVS) {
	HelpSubSysMap = map[string]HelpKVS{}
	for subSys, hkvs := range helpKVSMap {
		HelpSubSysMap[subSys] = hkvs
	}
}

// HelpDeprecatedSubSysMap - help for all deprecated sub-systems, that may be
// removed in the future.
var HelpDeprecatedSubSysMap map[string]HelpKV

// RegisterHelpDeprecatedSubSys - saves input help KVS for deprecated
// sub-systems globally. Should be called only once at init.
func RegisterHelpDeprecatedSubSys(helpDeprecatedKVMap map[string]HelpKV) {
	HelpDeprecatedSubSysMap = map[string]HelpKV{}
	for k, v := range helpDeprecatedKVMap {
		HelpDeprecatedSubSysMap[k] = v
	}
}

// KV - is a shorthand of each key value.
type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// KVS - is a shorthand for some wrapper functions
// to operate on list of key values.
type KVS []KV

// Empty - return if kv is empty
func (kvs KVS) Empty() bool {
	return len(kvs) == 0
}

// Clone - returns a copy of the KVS
func (kvs KVS) Clone() KVS {
	return append(make(KVS, 0, len(kvs)), kvs...)
}

// GetWithDefault - returns default value if key not set
func (kvs KVS) GetWithDefault(key string, defaultKVS KVS) string {
	v := kvs.Get(key)
	if len(v) == 0 {
		return defaultKVS.Get(key)
	}
	return v
}

// Keys returns the list of keys for the current KVS
func (kvs KVS) Keys() []string {
	keys := make([]string, len(kvs))
	var foundComment bool
	for i := range kvs {
		if kvs[i].Key == madmin.CommentKey {
			foundComment = true
		}
		keys[i] = kvs[i].Key
	}
	// Comment KV not found, add it explicitly.
	if !foundComment {
		keys = append(keys, madmin.CommentKey)
	}
	return keys
}

func (kvs KVS) String() string {
	var s strings.Builder
	for _, kv := range kvs {
		// Do not need to print if state is on
		if kv.Key == Enable && kv.Value == EnableOn {
			continue
		}
		s.WriteString(kv.Key)
		s.WriteString(KvSeparator)
		spc := madmin.HasSpace(kv.Value)
		if spc {
			s.WriteString(KvDoubleQuote)
		}
		s.WriteString(kv.Value)
		if spc {
			s.WriteString(KvDoubleQuote)
		}
		s.WriteString(KvSpaceSeparator)
	}
	return s.String()
}

// Merge environment values with on disk KVS, environment values overrides
// anything on the disk.
func Merge(cfgKVS map[string]KVS, envname string, defaultKVS KVS) map[string]KVS {
	newCfgKVS := make(map[string]KVS)
	for _, e := range env.List(envname) {
		tgt := strings.TrimPrefix(e, envname+Default)
		if tgt == envname {
			tgt = Default
		}
		newCfgKVS[tgt] = defaultKVS
	}
	for tgt, kv := range cfgKVS {
		newCfgKVS[tgt] = kv
	}
	return newCfgKVS
}

// Set sets a value, if not sets a default value.
func (kvs *KVS) Set(key, value string) {
	for i, kv := range *kvs {
		if kv.Key == key {
			(*kvs)[i] = KV{
				Key:   key,
				Value: value,
			}
			return
		}
	}
	*kvs = append(*kvs, KV{
		Key:   key,
		Value: value,
	})
}

// Get - returns the value of a key, if not found returns empty.
func (kvs KVS) Get(key string) string {
	v, ok := kvs.Lookup(key)
	if ok {
		return v
	}
	return ""
}

// Delete - deletes the key if present from the KV list.
func (kvs *KVS) Delete(key string) {
	for i, kv := range *kvs {
		if kv.Key == key {
			*kvs = append((*kvs)[:i], (*kvs)[i+1:]...)
			return
		}
	}
}

// Lookup - lookup a key in a list of KVS
func (kvs KVS) Lookup(key string) (string, bool) {
	for _, kv := range kvs {
		if kv.Key == key {
			return kv.Value, true
		}
	}
	return "", false
}

// Config - MinIO server config structure.
type Config map[string]map[string]KVS

// DelFrom - deletes all keys in the input reader.
func (c Config) DelFrom(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// Skip any empty lines, or comment like characters
		text := scanner.Text()
		if text == "" || strings.HasPrefix(text, KvComment) {
			continue
		}
		if err := c.DelKVS(text); err != nil {
			return err
		}
	}
	return scanner.Err()
}

// ReadConfig - read content from input and write into c.
// Returns whether all parameters were dynamic.
func (c Config) ReadConfig(r io.Reader) (dynOnly bool, err error) {
	var n int
	scanner := bufio.NewScanner(r)
	dynOnly = true
	for scanner.Scan() {
		// Skip any empty lines, or comment like characters
		text := scanner.Text()
		if text == "" || strings.HasPrefix(text, KvComment) {
			continue
		}
		dynamic, err := c.SetKVS(text, DefaultKVS)
		if err != nil {
			return false, err
		}
		dynOnly = dynOnly && dynamic
		n += len(text)
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return dynOnly, nil
}

// RedactSensitiveInfo - removes sensitive information
// like urls and credentials from the configuration
func (c Config) RedactSensitiveInfo() Config {
	nc := c.Clone()

	for configName, configVals := range nc {
		for _, helpKV := range HelpSubSysMap[configName] {
			if helpKV.Sensitive {
				for name, kvs := range configVals {
					for i := range kvs {
						if kvs[i].Key == helpKV.Key && len(kvs[i].Value) > 0 {
							kvs[i].Value = "*redacted*"
						}
					}
					configVals[name] = kvs
				}
			}
		}
	}

	// Remove the server credentials altogether
	nc.DelKVS(CredentialsSubSys)

	return nc
}

type configWriteTo struct {
	Config
	filterByKey string
}

// NewConfigWriteTo - returns a struct which
// allows for serializing the config/kv struct
// to a io.WriterTo
func NewConfigWriteTo(cfg Config, key string) io.WriterTo {
	return &configWriteTo{Config: cfg, filterByKey: key}
}

// WriteTo - implements io.WriterTo interface implementation for config.
func (c *configWriteTo) WriteTo(w io.Writer) (int64, error) {
	kvsTargets, err := c.GetKVS(c.filterByKey, DefaultKVS)
	if err != nil {
		return 0, err
	}
	var n int
	for _, target := range kvsTargets {
		m1, _ := w.Write([]byte(target.SubSystem))
		m2, _ := w.Write([]byte(KvSpaceSeparator))
		m3, _ := w.Write([]byte(target.KVS.String()))
		if len(kvsTargets) > 1 {
			m4, _ := w.Write([]byte(KvNewline))
			n += m1 + m2 + m3 + m4
		} else {
			n += m1 + m2 + m3
		}
	}
	return int64(n), nil
}

// Default KV configs for worm and region
var (
	DefaultCredentialKVS = KVS{
		KV{
			Key:   AccessKey,
			Value: auth.DefaultAccessKey,
		},
		KV{
			Key:   SecretKey,
			Value: auth.DefaultSecretKey,
		},
	}

	DefaultSiteKVS = KVS{
		KV{
			Key:   NameKey,
			Value: "",
		},
		KV{
			Key:   RegionKey,
			Value: "",
		},
	}

	DefaultRegionKVS = KVS{
		KV{
			Key:   RegionName,
			Value: "",
		},
	}
)

// LookupCreds - lookup credentials from config.
func LookupCreds(kv KVS) (auth.Credentials, error) {
	if err := CheckValidKeys(CredentialsSubSys, kv, DefaultCredentialKVS); err != nil {
		return auth.Credentials{}, err
	}
	accessKey := kv.Get(AccessKey)
	secretKey := kv.Get(SecretKey)
	if accessKey == "" || secretKey == "" {
		accessKey = auth.DefaultAccessKey
		secretKey = auth.DefaultSecretKey
	}
	return auth.CreateCredentials(accessKey, secretKey)
}

// Site - holds site info - name and region.
type Site struct {
	Name   string
	Region string
}

var validRegionRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-_-]+$")

// validSiteNameRegex - allows lowercase letters, digits and '-', starts with
// letter. At least 2 characters long.
var validSiteNameRegex = regexp.MustCompile("^[a-z][a-z0-9-]+$")

// LookupSite - get site related configuration. Loads configuration from legacy
// region sub-system as well.
func LookupSite(siteKV KVS, regionKV KVS) (s Site, err error) {
	if err = CheckValidKeys(SiteSubSys, siteKV, DefaultSiteKVS); err != nil {
		return
	}
	region := env.Get(EnvRegion, "")
	if region == "" {
		env.Get(EnvRegionName, "")
	}
	if region == "" {
		region = env.Get(EnvSiteRegion, siteKV.Get(RegionKey))
	}
	if region == "" {
		// No region config found in the site-subsystem. So lookup the legacy
		// region sub-system.
		if err = CheckValidKeys(RegionSubSys, regionKV, DefaultRegionKVS); err != nil {
			// An invalid key was found in the region sub-system.
			// Since the region sub-system cannot be (re)set as it
			// is legacy, we return an error to tell the user to
			// reset the region via the new command.
			err = Errorf("could not load region from legacy configuration as it was invalid - use 'mc admin config set myminio site region=myregion name=myname' to set a region and name (%v)", err)
			return
		}

		region = regionKV.Get(RegionName)
	}
	if region != "" {
		if !validRegionRegex.MatchString(region) {
			err = Errorf(
				"region '%s' is invalid, expected simple characters such as [us-east-1, myregion...]",
				region)
			return
		}
		s.Region = region
	}

	name := env.Get(EnvSiteName, siteKV.Get(NameKey))
	if name != "" {
		if !validSiteNameRegex.MatchString(name) {
			err = Errorf(
				"site name '%s' is invalid, expected simple characters such as [cal-rack0, myname...]",
				name)
			return
		}
		s.Name = name
	}
	return
}

// CheckValidKeys - checks if inputs KVS has the necessary keys,
// returns error if it find extra or superflous keys.
func CheckValidKeys(subSys string, kv KVS, validKVS KVS) error {
	nkv := KVS{}
	for _, kv := range kv {
		// Comment is a valid key, its also fully optional
		// ignore it since it is a valid key for all
		// sub-systems.
		if kv.Key == Comment {
			continue
		}
		if _, ok := validKVS.Lookup(kv.Key); !ok {
			nkv = append(nkv, kv)
		}
	}
	if len(nkv) > 0 {
		return Errorf(
			"found invalid keys (%s) for '%s' sub-system, use 'mc admin config reset myminio %s' to fix invalid keys", nkv.String(), subSys, subSys)
	}
	return nil
}

// LookupWorm - check if worm is enabled
func LookupWorm() (bool, error) {
	return ParseBool(env.Get(EnvWorm, EnableOff))
}

// Carries all the renamed sub-systems from their
// previously known names
var renamedSubsys = map[string]string{
	CrawlerSubSys: ScannerSubSys,
	// Add future sub-system renames
}

// Merge - merges a new config with all the
// missing values for default configs,
// returns a config.
func (c Config) Merge() Config {
	cp := New()
	for subSys, tgtKV := range c {
		for tgt := range tgtKV {
			ckvs := c[subSys][tgt]
			for _, kv := range cp[subSys][Default] {
				_, ok := c[subSys][tgt].Lookup(kv.Key)
				if !ok {
					ckvs.Set(kv.Key, kv.Value)
				}
			}
			if _, ok := cp[subSys]; !ok {
				rnSubSys, ok := renamedSubsys[subSys]
				if !ok {
					// A config subsystem was removed or server was downgraded.
					continue
				}
				// Copy over settings from previous sub-system
				// to newly renamed sub-system
				for _, kv := range cp[rnSubSys][Default] {
					_, ok := c[subSys][tgt].Lookup(kv.Key)
					if !ok {
						ckvs.Set(kv.Key, kv.Value)
					}
				}
				subSys = rnSubSys
			}
			cp[subSys][tgt] = ckvs
		}
	}
	return cp
}

// New - initialize a new server config.
func New() Config {
	srvCfg := make(Config)
	for _, k := range SubSystems.ToSlice() {
		srvCfg[k] = map[string]KVS{}
		srvCfg[k][Default] = DefaultKVS[k]
	}
	return srvCfg
}

// Target signifies an individual target
type Target struct {
	SubSystem string
	KVS       KVS
}

// Targets sub-system targets
type Targets []Target

// GetKVS - get kvs from specific subsystem.
func (c Config) GetKVS(s string, defaultKVS map[string]KVS) (Targets, error) {
	if len(s) == 0 {
		return nil, Errorf("input cannot be empty")
	}
	inputs := strings.Fields(s)
	if len(inputs) > 1 {
		return nil, Errorf("invalid number of arguments %s", s)
	}
	subSystemValue := strings.SplitN(inputs[0], SubSystemSeparator, 2)
	if len(subSystemValue) == 0 {
		return nil, Errorf("invalid number of arguments %s", s)
	}
	found := SubSystems.Contains(subSystemValue[0])
	if !found {
		// Check for sub-prefix only if the input value is only a
		// single value, this rejects invalid inputs if any.
		found = !SubSystems.FuncMatch(strings.HasPrefix, subSystemValue[0]).IsEmpty() && len(subSystemValue) == 1
	}
	if !found {
		return nil, Errorf("unknown sub-system %s", s)
	}

	targets := Targets{}
	subSysPrefix := subSystemValue[0]
	if len(subSystemValue) == 2 {
		if len(subSystemValue[1]) == 0 {
			return nil, Errorf("sub-system target '%s' cannot be empty", s)
		}
		kvs, ok := c[subSysPrefix][subSystemValue[1]]
		if !ok {
			return nil, Errorf("sub-system target '%s' doesn't exist", s)
		}
		for _, kv := range defaultKVS[subSysPrefix] {
			_, ok = kvs.Lookup(kv.Key)
			if !ok {
				kvs.Set(kv.Key, kv.Value)
			}
		}
		targets = append(targets, Target{
			SubSystem: inputs[0],
			KVS:       kvs,
		})
	} else {
		// Use help for sub-system to preserve the order. Add deprecated
		// keys at the end (in some order).
		kvsOrder := append([]HelpKV{}, HelpSubSysMap[""]...)
		for _, v := range HelpDeprecatedSubSysMap {
			kvsOrder = append(kvsOrder, v)
		}

		for _, hkv := range kvsOrder {
			if !strings.HasPrefix(hkv.Key, subSysPrefix) {
				continue
			}
			if c[hkv.Key][Default].Empty() {
				targets = append(targets, Target{
					SubSystem: hkv.Key,
					KVS:       defaultKVS[hkv.Key],
				})
			}
			for k, kvs := range c[hkv.Key] {
				for _, dkv := range defaultKVS[hkv.Key] {
					_, ok := kvs.Lookup(dkv.Key)
					if !ok {
						kvs.Set(dkv.Key, dkv.Value)
					}
				}
				if k != Default {
					targets = append(targets, Target{
						SubSystem: hkv.Key + SubSystemSeparator + k,
						KVS:       kvs,
					})
				} else {
					targets = append(targets, Target{
						SubSystem: hkv.Key,
						KVS:       kvs,
					})
				}
			}
		}
	}
	return targets, nil
}

// DelKVS - delete a specific key.
func (c Config) DelKVS(s string) error {
	if len(s) == 0 {
		return Errorf("input arguments cannot be empty")
	}
	inputs := strings.Fields(s)
	if len(inputs) > 1 {
		return Errorf("invalid number of arguments %s", s)
	}
	subSystemValue := strings.SplitN(inputs[0], SubSystemSeparator, 2)
	if len(subSystemValue) == 0 {
		return Errorf("invalid number of arguments %s", s)
	}
	if !SubSystems.Contains(subSystemValue[0]) {
		// Unknown sub-system found try to remove it anyways.
		delete(c, subSystemValue[0])
		return nil
	}
	tgt := Default
	subSys := subSystemValue[0]
	if len(subSystemValue) == 2 {
		if len(subSystemValue[1]) == 0 {
			return Errorf("sub-system target '%s' cannot be empty", s)
		}
		tgt = subSystemValue[1]
	}
	_, ok := c[subSys][tgt]
	if !ok {
		return Errorf("sub-system %s already deleted", s)
	}
	delete(c[subSys], tgt)
	return nil
}

// Clone - clones a config map entirely.
func (c Config) Clone() Config {
	cp := New()
	for subSys, tgtKV := range c {
		cp[subSys] = make(map[string]KVS)
		for tgt, kv := range tgtKV {
			cp[subSys][tgt] = append(cp[subSys][tgt], kv...)
		}
	}
	return cp
}

// GetSubSys - extracts subssystem info from given config string
func GetSubSys(s string) (subSys string, inputs []string, tgt string, e error) {
	tgt = Default
	if len(s) == 0 {
		return subSys, inputs, tgt, Errorf("input arguments cannot be empty")
	}
	inputs = strings.SplitN(s, KvSpaceSeparator, 2)

	subSystemValue := strings.SplitN(inputs[0], SubSystemSeparator, 2)
	subSys = subSystemValue[0]
	if !SubSystems.Contains(subSys) {
		return subSys, inputs, tgt, Errorf("unknown sub-system %s", s)
	}

	if len(inputs) == 1 {
		return subSys, inputs, tgt, nil
	}

	if SubSystemsSingleTargets.Contains(subSystemValue[0]) && len(subSystemValue) == 2 {
		return subSys, inputs, tgt, Errorf("sub-system '%s' only supports single target", subSystemValue[0])
	}

	if len(subSystemValue) == 2 {
		tgt = subSystemValue[1]
	}

	return subSys, inputs, tgt, e
}

// SetKVS - set specific key values per sub-system.
func (c Config) SetKVS(s string, defaultKVS map[string]KVS) (dynamic bool, err error) {
	subSys, inputs, tgt, err := GetSubSys(s)
	if err != nil {
		return false, err
	}

	dynamic = SubSystemsDynamic.Contains(subSys)

	fields := madmin.KvFields(inputs[1], defaultKVS[subSys].Keys())
	if len(fields) == 0 {
		return false, Errorf("sub-system '%s' cannot have empty keys", subSys)
	}

	kvs := KVS{}
	var prevK string
	for _, v := range fields {
		kv := strings.SplitN(v, KvSeparator, 2)
		if len(kv) == 0 {
			continue
		}
		if len(kv) == 1 && prevK != "" {
			value := strings.Join([]string{
				kvs.Get(prevK),
				madmin.SanitizeValue(kv[0]),
			}, KvSpaceSeparator)
			kvs.Set(prevK, value)
			continue
		}
		if len(kv) == 2 {
			prevK = kv[0]
			kvs.Set(prevK, madmin.SanitizeValue(kv[1]))
			continue
		}
		return false, Errorf("key '%s', cannot have empty value", kv[0])
	}

	_, ok := kvs.Lookup(Enable)
	// Check if state is required
	_, enableRequired := defaultKVS[subSys].Lookup(Enable)
	if !ok && enableRequired {
		// implicit state "on" if not specified.
		kvs.Set(Enable, EnableOn)
	}

	var currKVS KVS
	ck, ok := c[subSys][tgt]
	if !ok {
		currKVS = defaultKVS[subSys].Clone()
	} else {
		currKVS = ck.Clone()
		for _, kv := range defaultKVS[subSys] {
			if _, ok = currKVS.Lookup(kv.Key); !ok {
				currKVS.Set(kv.Key, kv.Value)
			}
		}
	}

	for _, kv := range kvs {
		if kv.Key == Comment {
			// Skip comment and add it later.
			continue
		}
		currKVS.Set(kv.Key, kv.Value)
	}

	v, ok := kvs.Lookup(Comment)
	if ok {
		currKVS.Set(Comment, v)
	}

	hkvs := HelpSubSysMap[subSys]
	for _, hkv := range hkvs {
		var enabled bool
		if enableRequired {
			enabled = currKVS.Get(Enable) == EnableOn
		} else {
			// when enable arg is not required
			// then it is implicit on for the sub-system.
			enabled = true
		}
		v, _ := currKVS.Lookup(hkv.Key)
		if v == "" && !hkv.Optional && enabled {
			// Return error only if the
			// key is enabled, for state=off
			// let it be empty.
			return false, Errorf(
				"'%s' is not optional for '%s' sub-system, please check '%s' documentation",
				hkv.Key, subSys, subSys)
		}
	}
	c[subSys][tgt] = currKVS
	return dynamic, nil
}

// CheckValidKeys - checks if the config parameters for the given subsystem and
// target are valid. It checks both the configuration store as well as
// environment variables.
func (c Config) CheckValidKeys(subSys string, deprecatedKeys []string) error {
	defKVS, ok := DefaultKVS[subSys]
	if !ok {
		return fmt.Errorf("Subsystem %s does not exist", subSys)
	}

	// Make a list of valid keys for the subsystem including the `comment`
	// key.
	validKeys := make([]string, 0, len(defKVS)+1)
	for _, param := range defKVS {
		validKeys = append(validKeys, param.Key)
	}
	validKeys = append(validKeys, Comment)

	subSysEnvVars := env.List(fmt.Sprintf("%s%s", EnvPrefix, strings.ToUpper(subSys)))

	// Set of env vars for the sub-system to validate.
	candidates := set.CreateStringSet(subSysEnvVars...)

	// Remove all default target env vars from the candidates set (as they
	// are valid).
	for _, param := range validKeys {
		paramEnvName := getEnvVarName(subSys, Default, param)
		candidates.Remove(paramEnvName)
	}

	isSingleTarget := SubSystemsSingleTargets.Contains(subSys)
	if isSingleTarget && len(candidates) > 0 {
		return fmt.Errorf("The following environment variables are unknown: %s",
			strings.Join(candidates.ToSlice(), ", "))
	}

	if !isSingleTarget {
		// Validate other env vars for all targets.
		envVars := candidates.ToSlice()
		for _, envVar := range envVars {
			for _, param := range validKeys {
				pEnvName := getEnvVarName(subSys, Default, param) + Default
				if len(envVar) > len(pEnvName) && strings.HasPrefix(envVar, pEnvName) {
					// This envVar is valid - it has a
					// non-empty target.
					candidates.Remove(envVar)
				}
			}
		}

		// Whatever remains are invalid env vars - return an error.
		if len(candidates) > 0 {
			return fmt.Errorf("The following environment variables are unknown: %s",
				strings.Join(candidates.ToSlice(), ", "))
		}
	}

	validKeysSet := set.CreateStringSet(validKeys...)
	validKeysSet = validKeysSet.Difference(set.CreateStringSet(deprecatedKeys...))
	kvsMap := c[subSys]
	for tgt, kvs := range kvsMap {
		invalidKV := KVS{}
		for _, kv := range kvs {
			if !validKeysSet.Contains(kv.Key) {
				invalidKV = append(invalidKV, kv)
			}
		}
		if len(invalidKV) > 0 {
			return Errorf(
				"found invalid keys (%s) for '%s:%s' sub-system, use 'mc admin config reset myminio %s:%s' to fix invalid keys",
				invalidKV.String(), subSys, tgt, subSys, tgt)
		}
	}
	return nil
}

// GetAvailableTargets - returns a list of targets configured for the given
// subsystem (whether they are enabled or not). A target could be configured via
// environment variables or via the configuration store. The default target is
// `_` and is always returned.
func (c Config) GetAvailableTargets(subSys string) ([]string, error) {
	if SubSystemsSingleTargets.Contains(subSys) {
		return []string{Default}, nil
	}

	defKVS, ok := DefaultKVS[subSys]
	if !ok {
		return nil, fmt.Errorf("Subsystem %s does not exist", subSys)
	}

	kvsMap := c[subSys]
	s := set.NewStringSet()

	// Add all targets that are configured in the config store.
	for k := range kvsMap {
		s.Add(k)
	}

	// Add targets that are configured via environment variables.
	for _, param := range defKVS {
		envVarPrefix := getEnvVarName(subSys, Default, param.Key) + Default
		envsWithPrefix := env.List(envVarPrefix)
		for _, k := range envsWithPrefix {
			tgtName := strings.TrimPrefix(k, envVarPrefix)
			if tgtName != "" {
				s.Add(tgtName)
			}
		}
	}

	return s.ToSlice(), nil
}

func getEnvVarName(subSys, target, param string) string {
	if target == Default {
		return fmt.Sprintf("%s%s_%s", EnvPrefix, strings.ToUpper(subSys), strings.ToUpper(param))
	}

	return fmt.Sprintf("%s%s_%s_%s", EnvPrefix, strings.ToUpper(subSys), strings.ToUpper(param), target)
}

var resolvableSubsystems = set.CreateStringSet(IdentityOpenIDSubSys)

// ValueSource represents the source of a config parameter value.
type ValueSource uint8

// Constants for ValueSource
const (
	ValueSourceAbsent ValueSource = iota // this is an error case
	ValueSourceDef
	ValueSourceCfg
	ValueSourceEnv
)

// ResolveConfigParam returns the effective value of a configuration parameter,
// within a subsystem and subsystem target. The effective value is, in order of
// decreasing precedence:
//
// 1. the value of the corresponding environment variable if set,
// 2. the value of the parameter in the config store if set,
// 3. the default value,
//
// This function only works for a subset of sub-systems, others return
// `ValueSourceAbsent`. FIXME: some parameters have custom environment
// variables for which support needs to be added.
func (c Config) ResolveConfigParam(subSys, target, cfgParam string) (value string, cs ValueSource) {
	// cs = ValueSourceAbsent initially as it is iota by default.

	// Initially only support OpenID
	if !resolvableSubsystems.Contains(subSys) {
		return
	}

	// Check if config param requested is valid.
	defKVS, ok := DefaultKVS[subSys]
	if !ok {
		return
	}

	defValue, isFound := defKVS.Lookup(cfgParam)
	if !isFound {
		return
	}

	if target == "" {
		target = Default
	}

	envVar := getEnvVarName(subSys, target, cfgParam)

	// Lookup Env var.
	value = env.Get(envVar, "")
	if value != "" {
		cs = ValueSourceEnv
		return
	}

	// Lookup config store.
	if subSysStore, ok := c[subSys]; ok {
		if kvs, ok2 := subSysStore[target]; ok2 {
			var ok3 bool
			value, ok3 = kvs.Lookup(cfgParam)
			if ok3 {
				cs = ValueSourceCfg
				return
			}
		}
	}

	// Return the default value.
	value = defValue
	cs = ValueSourceDef
	return
}
