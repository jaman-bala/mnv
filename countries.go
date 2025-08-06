package mnv

// CountryPhoneCodes содержит информацию о телефонных кодах стран
var CountryPhoneCodes = map[string]PhoneCodeInfo{
	"kg": {Prefix: "+996", Pattern: `^\+996[0-9]{9}$`, MinLength: 9, MaxLength: 9},
	"ru": {Prefix: "+7", Pattern: `^\+7[0-9]{10}$`, MinLength: 10, MaxLength: 10},
	"kz": {Prefix: "+7", Pattern: `^\+7[0-9]{10}$`, MinLength: 10, MaxLength: 10},
	"uz": {Prefix: "+998", Pattern: `^\+998[0-9]{9}$`, MinLength: 9, MaxLength: 9},
	"us": {Prefix: "+1", Pattern: `^\+1[0-9]{10}$`, MinLength: 10, MaxLength: 10},
	"uk": {Prefix: "+44", Pattern: `^\+44[0-9]{10,11}$`, MinLength: 10, MaxLength: 11},
	"de": {Prefix: "+49", Pattern: `^\+49[0-9]{10,12}$`, MinLength: 10, MaxLength: 12},
	"fr": {Prefix: "+33", Pattern: `^\+33[0-9]{9,10}$`, MinLength: 9, MaxLength: 10},
	"it": {Prefix: "+39", Pattern: `^\+39[0-9]{9,11}$`, MinLength: 9, MaxLength: 11},
	"tr": {Prefix: "+90", Pattern: `^\+90[0-9]{10}$`, MinLength: 10, MaxLength: 10},
	"am": {Prefix: "+374", Pattern: `^\+374[0-9]{8}$`, MinLength: 8, MaxLength: 8},
	"az": {Prefix: "+994", Pattern: `^\+994[0-9]{9}$`, MinLength: 9, MaxLength: 9},
	"tj": {Prefix: "+992", Pattern: `^\+992[0-9]{9}$`, MinLength: 9, MaxLength: 9},
	"tm": {Prefix: "+993", Pattern: `^\+993[0-9]{8}$`, MinLength: 8, MaxLength: 8},
	"ua": {Prefix: "+380", Pattern: `^\+380[0-9]{9}$`, MinLength: 9, MaxLength: 9},
	"by": {Prefix: "+375", Pattern: `^\+375[0-9]{9}$`, MinLength: 9, MaxLength: 9},
	"cn": {Prefix: "+86", Pattern: `^\+86[0-9]{11}$`, MinLength: 11, MaxLength: 11},
	"in": {Prefix: "+91", Pattern: `^\+91[0-9]{10}$`, MinLength: 10, MaxLength: 10},
	"jp": {Prefix: "+81", Pattern: `^\+81[0-9]{10,11}$`, MinLength: 10, MaxLength: 11},
}
