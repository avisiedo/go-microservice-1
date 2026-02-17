//go:build ignore

// build ignore
package config

// func TestSetClowderConfiguration(t *testing.T) {
// 	// Panic for v = nil
// 	assert.PanicsWithValue(t, "'v' is nil", func() {
// 		setClowderConfiguration(nil, nil)
// 	})

// 	// Panic for cfg = nil
// 	v := viper.New()
// 	assert.PanicsWithValue(t, "'clowderConfig' is nil", func() {
// 		setClowderConfiguration(v, nil)
// 	})

// 	// Loading empty clowder.AppConfig
// 	cfg := clowder.AppConfig{}
// 	assert.NotPanics(t, func() {
// 		setClowderConfiguration(v, &cfg)
// 	})
// 	assert.Nil(t, v.Get("database.host"))
// 	assert.Nil(t, v.Get("database.port"))
// 	assert.Nil(t, v.Get("database.user"))
// 	assert.Nil(t, v.Get("database.password"))
// 	assert.Nil(t, v.Get("database.name"))
// 	assert.Nil(t, v.Get("database.ca_cert_path"))

// 	// Load RDSACert with wrong path
// 	cfg.Database = &clowder.DatabaseConfig{}
// 	cfg.Database.RdsCa = pointy.String("/tmp/itdoesnotexist.pem")
// 	assert.NotPanics(t, func() {
// 		setClowderConfiguration(v, &cfg)
// 	})

// 	// Load RDSACert nil
// 	cfg.Database.RdsCa = nil
// 	assert.NotPanics(t, func() {
// 		setClowderConfiguration(v, &cfg)
// 	})

// 	// Load database data (but RdsCa)
// 	cfg.Database = &clowder.DatabaseConfig{
// 		Hostname: "testhost",
// 		Port:     5432,
// 		Username: "testuser",
// 		Password: "testpassword",
// 		Name:     "testname",
// 	}
// 	assert.NotPanics(t, func() {
// 		setClowderConfiguration(v, &cfg)
// 	})
// 	assert.Equal(t, "testhost", v.Get("database.host"))
// 	assert.Equal(t, 5432, v.Get("database.port"))
// 	assert.Equal(t, "testuser", v.Get("database.user"))
// 	assert.Equal(t, "testpassword", v.Get("database.password"))
// 	assert.Equal(t, "testname", v.Get("database.name"))
// 	assert.NotNil(t, v.Get("database.ca_cert_path"))

// 	// TODO Add test to cover RdsCa flow

// 	// Load cloudwatch data
// 	cfg.Logging.Cloudwatch = &clowder.CloudWatchConfig{
// 		Region:          "testregion",
// 		LogGroup:        "testgroup",
// 		SecretAccessKey: "testsecret",
// 		AccessKeyId:     "testaccesskeyid",
// 	}
// 	assert.NotPanics(t, func() {
// 		setClowderConfiguration(v, &cfg)
// 	})
// 	assert.Equal(t, "testregion", v.Get("cloudwatch.region"))
// 	assert.Equal(t, "testgroup", v.Get("cloudwatch.group"))
// 	assert.Equal(t, "testsecret", v.Get("cloudwatch.secret"))
// 	assert.Equal(t, "testaccesskeyid", v.Get("cloudwatch.key"))
// }
