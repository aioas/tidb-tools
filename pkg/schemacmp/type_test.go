package schemacmp_test

import (
	. "github.com/pingcap/tidb-tools/pkg/schemacmp"

	. "github.com/pingcap/check"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/types"
)

type typeSchema struct{}

var _ = Suite(&typeSchema{})

const binary = "binary"

var (
	// INT
	typeInt = &types.FieldType{
		Tp:      mysql.TypeLong,
		Flag:    0,
		Flen:    11,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// INT UNSIGNED NOT NULL
	typeIntUnsignedNotNull = &types.FieldType{
		Tp:      mysql.TypeLong,
		Flag:    mysql.NoDefaultValueFlag | mysql.UnsignedFlag | mysql.NotNullFlag,
		Flen:    10,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// INT AUTO_INCREMENT UNIQUE
	typeIntAutoIncrementUnique = &types.FieldType{
		Tp:      mysql.TypeLong,
		Flag:    mysql.AutoIncrementFlag | mysql.UniqueKeyFlag,
		Flen:    11,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// INT(22)
	typeInt22 = &types.FieldType{
		Tp:      mysql.TypeLong,
		Flag:    0,
		Flen:    22,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// BIT(4)
	typeBit4 = &types.FieldType{
		Tp:      mysql.TypeBit,
		Flag:    mysql.UnsignedFlag,
		Flen:    4,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// BIGINT(22) ZEROFILL
	typeBigInt22ZeroFill = &types.FieldType{
		Tp:      mysql.TypeLonglong,
		Flag:    mysql.ZerofillFlag | mysql.UnsignedFlag,
		Flen:    22,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// DECIMAL(16, 8) DEFAULT 2.5
	typeDecimal16_8 = &types.FieldType{
		Tp:      mysql.TypeNewDecimal,
		Flag:    0,
		Flen:    16,
		Decimal: 8,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// DECIMAL
	typeDecimal = &types.FieldType{
		Tp:      mysql.TypeNewDecimal,
		Flag:    0,
		Flen:    11,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// DATE
	typeDate = &types.FieldType{
		Tp:      mysql.TypeDate,
		Flag:    mysql.BinaryFlag,
		Flen:    10,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// DATETIME(3)
	typeDateTime3 = &types.FieldType{
		Tp:      mysql.TypeDatetime,
		Flag:    mysql.BinaryFlag,
		Flen:    23,
		Decimal: 3,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// TIMESTAMP
	typeTimestamp = &types.FieldType{
		Tp:      mysql.TypeTimestamp,
		Flag:    mysql.BinaryFlag,
		Flen:    19,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// TIME(6)
	typeTime6 = &types.FieldType{
		Tp:      mysql.TypeDuration,
		Flag:    mysql.BinaryFlag,
		Flen:    17,
		Decimal: 6,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// YEAR(4)
	typeYear4 = &types.FieldType{
		Tp:      mysql.TypeYear,
		Flag:    mysql.ZerofillFlag | mysql.UnsignedFlag,
		Flen:    4,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// CHAR(123)
	typeChar123 = &types.FieldType{
		Tp:      mysql.TypeString,
		Flag:    0,
		Flen:    123,
		Decimal: 0,
		Charset: mysql.UTF8MB4Charset,
		Collate: mysql.UTF8MB4DefaultCollation,
		Elems:   nil,
	}

	// VARCHAR(65432) CHARSET ascii
	typeVarchar65432CharsetASCII = &types.FieldType{
		Tp:      mysql.TypeVarchar,
		Flag:    0,
		Flen:    65432,
		Decimal: 0,
		Charset: "ascii",
		Collate: "ascii_bin",
		Elems:   nil,
	}

	// BINARY(69)
	typeBinary69 = &types.FieldType{
		Tp:      mysql.TypeString,
		Flag:    mysql.BinaryFlag,
		Flen:    69,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// VARBINARY(420)
	typeVarBinary420 = &types.FieldType{
		Tp:      mysql.TypeVarchar,
		Flag:    mysql.BinaryFlag,
		Flen:    420,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// LONGBLOB
	typeLongBlob = &types.FieldType{
		Tp:      mysql.TypeLongBlob,
		Flag:    mysql.BinaryFlag,
		Flen:    0xffffffff,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}

	// MEDIUMTEXT
	typeMediumText = &types.FieldType{
		Tp:      mysql.TypeMediumBlob,
		Flag:    0,
		Flen:    0xffffff,
		Decimal: 0,
		Charset: mysql.UTF8MB4Charset,
		Collate: mysql.UTF8MB4DefaultCollation,
		Elems:   nil,
	}

	// ENUM('tidb', 'tikv', 'tiflash', 'golang', 'rust')
	typeEnum5 = &types.FieldType{
		Tp:      mysql.TypeEnum,
		Flag:    0,
		Flen:    types.UnspecifiedLength,
		Decimal: 0,
		Charset: mysql.UTF8MB4Charset,
		Collate: mysql.UTF8MB4DefaultCollation,
		Elems:   []string{"tidb", "tikv", "tiflash", "golang", "rust"},
	}

	// ENUM('tidb', 'tikv')
	typeEnum2 = &types.FieldType{
		Tp:      mysql.TypeEnum,
		Flag:    0,
		Flen:    types.UnspecifiedLength,
		Decimal: 0,
		Charset: mysql.UTF8MB4Charset,
		Collate: mysql.UTF8MB4DefaultCollation,
		Elems:   []string{"tidb", "tikv"},
	}

	// SET('tidb', 'tikv', 'tiflash', 'golang', 'rust')
	typeSet5 = &types.FieldType{
		Tp:      mysql.TypeSet,
		Flag:    0,
		Flen:    types.UnspecifiedLength,
		Decimal: 0,
		Charset: mysql.UTF8MB4Charset,
		Collate: mysql.UTF8MB4DefaultCollation,
		Elems:   []string{"tidb", "tikv", "tiflash", "golang", "rust"},
	}

	// ENUM('tidb', 'tikv')
	typeSet2 = &types.FieldType{
		Tp:      mysql.TypeSet,
		Flag:    0,
		Flen:    types.UnspecifiedLength,
		Decimal: 0,
		Charset: mysql.UTF8MB4Charset,
		Collate: mysql.UTF8MB4DefaultCollation,
		Elems:   []string{"tidb", "tikv"},
	}

	// JSON
	typeJSON = &types.FieldType{
		Tp:      mysql.TypeJSON,
		Flag:    mysql.BinaryFlag,
		Flen:    0xffffffff,
		Decimal: 0,
		Charset: binary,
		Collate: binary,
		Elems:   nil,
	}
)

func (*typeSchema) TestTypeUnwrap(c *C) {
	testCases := []*types.FieldType{
		typeInt,
		typeIntUnsignedNotNull,
		typeIntAutoIncrementUnique,
		typeInt22,
		typeBit4,
		typeBigInt22ZeroFill,
		typeDecimal16_8,
		typeDecimal,
		typeDate,
		typeDateTime3,
		typeTimestamp,
		typeTime6,
		typeYear4,
		typeChar123,
		typeVarchar65432CharsetASCII,
		typeBinary69,
		typeVarBinary420,
		typeLongBlob,
		typeMediumText,
		typeEnum5,
		typeEnum2,
		typeSet5,
		typeSet2,
		typeJSON,
	}

	for _, tc := range testCases {
		assert := func(expected interface{}, checker Checker, args ...interface{}) {
			c.Assert(expected, checker, append(args, Commentf("tc = %s", tc))...)
		}
		t := Type(tc)
		assert(t.Unwrap(), DeepEquals, tc)
	}
}
