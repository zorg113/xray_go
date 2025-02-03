package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	IncorrectStruct1 struct {
		Data string `validate:"min:12"`
	}
	IncorrectStruct2 struct {
		Data string `validate:"hz:12"`
	}

	IncorrectStruct3 struct {
		Data string `validate:""`
	}

	IncorrectStruct4 struct {
		Data string `validate:"regexp:^\\w+@\\w+*\\.\\W+$"`
	}

	TestStruct struct {
		Data1 int         `validate:"min:12"`
		Data2 TestStruct1 `validate:"nested"`
	}

	TestStruct1 struct {
		Data1_2 TestStruct2 `validate:"nested"`
		Data2_2 int         `validate:"max:10"`
	}

	TestStruct2 struct {
		Data1_3 int              `validate:"in:1,2"`
		Data2_3 string           `validate:"in:3,4"`
		Data3_3 IncorrectStruct4 // no validation
		Data4_3 IncorrectStruct2 // no validateion
	}
	TestStruct3 struct {
		Data1_3 int              `validate:"in:2"`
		Data2_3 IncorrectStruct1 `validate:"nested"`
	}
)

var tests = []struct {
	in          interface{}
	expectedErr error
}{
	{
		in: User{
			ID:     "896875896875896875896875896875896875", // string `json:"id" validate:"len:36"`
			Age:    49,                                     // int `validate:"min:18|max:50"`
			Email:  "jon@deer.uk",                          // string `validate`
			Role:   "admin",                                // string `validate:"in:admin,stuff"`
			Phones: nil,                                    // []string `validate:"len:11"`
		},
		expectedErr: nil,
	},
	{
		in: User{
			ID:     "8968",        // string `json:"id" validate:"len:36"`
			Age:    49,            // int `validate:"min:18|max:50"`
			Email:  "jon@deer.uk", // string `validate`
			Role:   "admin",
			Phones: []string{"89009999999", "89002999999"},
			meta:   json.RawMessage{},
		},
		expectedErr: ValidationErrors{{Field: "ID", Err: ErrLen}},
	},
	{
		in: User{
			ID:     "896875896875896875896875896875896875", // string `json:"id" validate:"len:36"`
			Age:    15,                                     // int `validate:"min:18|max:50"`
			Email:  "jon@deer.uk",                          // string `validate`
			Role:   "admin",
			Phones: []string{"89009999999", "89002999999"},
			meta:   json.RawMessage{},
		},
		expectedErr: ValidationErrors{{Field: "Age", Err: ErrMin}},
	},
	{
		in: User{
			ID:     "896875896875896875896875896875896875", // string `json:"id" validate:"len:36"`
			Age:    150,                                    // int `validate:"min:18|max:50"`
			Email:  "jon@deer.uk",                          // string `validate`
			Role:   "admin",
			Phones: []string{"89009999999", "89002999999"},
			meta:   json.RawMessage{},
		},
		expectedErr: ValidationErrors{{Field: "Age", Err: ErrMax}},
	},
	{
		in: User{
			ID:     "896875896875896875896875896875896875", // string `json:"id" validate:"len:36"`
			Age:    25,                                     // int `validate:"min:18|max:50"`
			Email:  "jondeer.uk",                           // string `validate`
			Role:   "stuff",
			Phones: []string{"89009999999", "89002999999"},
			meta:   json.RawMessage{},
		},
		expectedErr: ValidationErrors{{Field: "Email", Err: ErrRegExp}},
	},
	{
		in: User{
			ID:     "896875896875896875896875896875896875", // string `json:"id" validate:"len:36"`
			Age:    49,                                     // int `validate:"min:18|max:50"`
			Email:  "jon@deer.uk",                          // string `validate`
			Role:   "manager",                              // string `validate:"in:admin,stuff"`
			Phones: nil,                                    // []string `validate:"len:11"`
			meta:   json.RawMessage{},
		},
		expectedErr: ValidationErrors{{Field: "Role", Err: ErrIn}},
	},
	{
		in: User{
			ID:     "896875896875896875896875896875896875", // string `json:"id" validate:"len:36"`
			Age:    49,                                     // int `validate:"min:18|max:50"`
			Email:  "jon@deer.uk",                          // string `validate`
			Role:   "admin",                                // string `validate:"in:admin,stuff"`
			Phones: []string{"890099999", "890029999949"},  // []string `validate:"len:11"`
			meta:   json.RawMessage{},
		},
		expectedErr: ValidationErrors{{Field: "Phones", Err: ErrLen}},
	},
	{
		in: User{
			ID:     "8906875",                             // string `json:"id" validate:"len:36"`
			Age:    4,                                     // int `validate:"min:18|max:50"`
			Email:  "jon@deer",                            // string `validate`
			Role:   "manager",                             // string `validate:"in:admin,stuff"`
			Phones: []string{"890099999", "890029999949"}, // []string `validate:"len:11"`
			meta:   json.RawMessage{},
		},
		expectedErr: ValidationErrors{
			{Field: "ID", Err: ErrLen},
			{Field: "Age", Err: ErrMin},
			{Field: "Email", Err: ErrRegExp},
			{Field: "Role", Err: ErrIn},
			{Field: "Phones", Err: ErrLen},
		},
	},
	{
		in:          App{Version: "22334"},
		expectedErr: nil,
	},
	{
		in:          App{Version: "12345"},
		expectedErr: nil,
	},
	{
		in:          App{Version: "123456"},
		expectedErr: ValidationErrors{{Field: "Version", Err: ErrLen}},
	},
	{
		in:          App{Version: ""},
		expectedErr: ValidationErrors{{Field: "Version", Err: ErrLen}},
	},

	{
		in:          Token{},
		expectedErr: nil,
	},

	{
		in: Response{
			Code: 200, // int    `validate:"in:200,404,500"`
			Body: "",  // string `json:"omitempty"`
		},
		expectedErr: nil,
	},
	{
		in: Response{
			Code: 0,  // int    `validate:"in:200,404,500"`
			Body: "", // string `json:"omitempty"`
		},
		expectedErr: ValidationErrors{{Field: "Code", Err: ErrIn}},
	},

	{in: IncorrectStruct1{}, expectedErr: ErrTag},
	{in: IncorrectStruct2{}, expectedErr: ErrTag},
	{in: IncorrectStruct3{}, expectedErr: ErrTag},
	{in: IncorrectStruct4{}, expectedErr: ErrTag},
	{
		in: TestStruct{
			Data1: 1, // err
			Data2: TestStruct1{
				Data1_2: TestStruct2{
					Data1_3: 1,   // ok
					Data2_3: "5", // err
				},
				Data2_2: 11, // err
			},
		},
		expectedErr: ValidationErrors{
			{Field: "Data1", Err: ErrMin},
			{Field: "Data2_3", Err: ErrIn},
			{Field: "Data2_2", Err: ErrMax},
		},
	},

	{in: TestStruct3{}, expectedErr: ErrTag},
	{in: []string{}, expectedErr: ErrNotAStruct},
}

func TestValidate(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.Nil(t, err)
				return
			}
			var expErrs ValidationErrors
			if errors.As(tt.expectedErr, &expErrs) {
				var resErrs ValidationErrors
				require.ErrorAs(t, err, &resErrs)
				require.Equal(t, len(expErrs), len(resErrs))
				for i := range expErrs {
					require.ErrorIs(t, expErrs[i].Err, resErrs[i].Err)
					require.Equal(t, expErrs[i].Field, resErrs[i].Field)
				}
			} else {
				require.ErrorIs(t, tt.expectedErr, err)
			}
		})
	}
}
