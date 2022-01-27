# GORM ClickHouse Driver

Clickhouse support for GORM

[![test status](https://github.com/go-gorm/clickhouse/workflows/tests/badge.svg?branch=master "test status")](https://github.com/go-gorm/clickhouse/actions)

## Quick Start

You can simply test your connection to your database with the following:

```go
import (
  "gorm.io/driver/clickhouse"
  "gorm.io/gorm"
)

func main() {
  dsn := "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20"
  db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})

  // Auto Migrate
  db.AutoMigrate(&User{})
  // Set table options
  db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&User{})

  // Insert
  db.Create(&user)

  // Select
  db.Find(&user, "id = ?", 10)

  // Batch Insert
  var users = []User{user1, user2, user3}
  db.Create(&users)
  // ...
}
```

## Advanced Configuration

```go
import (
  "gorm.io/driver/clickhouse"
  "gorm.io/gorm"
)

// refer to https://github.com/ClickHouse/clickhouse-go
var dsn = "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20"

func main() {
  db, err := gorm.Open(clickhouse.New(click.Config{
    DSN: dsn,
    Conn: conn,                       // initialize with existing database conn
    DisableDatetimePrecision: true,   // disable datetime64 precision, not supported before clickhouse 20.4
    DontSupportRenameColumn: true,    // rename column not supported before clickhouse 20.4
    SkipInitializeWithVersion: false, // smart configure based on used version
    DefaultGranularity: 3,            // 1 granule = 8192 rows
    DefaultCompression: "LZ4",        // default compression algorithm. LZ4 is lossless
    DefaultIndexType: "minmax",       // index stores extremes of the expression
    DefaultTableEngineOpts: "ENGINE=MergeTree() ORDER BY tuple()",
  }), &gorm.Config{})
}
```

Checkout [https://gorm.io](https://gorm.io) for details.
