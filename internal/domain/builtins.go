package domain

import "strings"

// BuiltinSymbolSet provides O(1) zero-memory set lookups for language built-ins and dunder methods.
var BuiltinSymbolSet = map[string]map[string]struct{}{
	"py": {
		// Special dunder methods & attributes
		"__init__": {}, "__str__": {}, "__repr__": {}, "__len__": {}, "__eq__": {},
		"__ne__": {}, "__getitem__": {}, "__setitem__": {}, "__delitem__": {},
		"__iter__": {}, "__next__": {}, "__call__": {}, "__enter__": {}, "__exit__": {},
		"__main__": {}, "__name__": {}, "__doc__": {}, "__file__": {}, "__all__": {},
		"__path__": {}, "__package__": {}, "__dict__": {}, "__slots__": {},
		// Python Built-in Functions
		"print": {}, "len": {}, "range": {}, "enumerate": {}, "open": {}, "input": {},
		"str": {}, "int": {}, "float": {}, "bool": {}, "list": {}, "dict": {}, "set": {},
		"tuple": {}, "super": {}, "isinstance": {}, "issubclass": {}, "type": {},
		"abs": {}, "all": {}, "any": {}, "bin": {}, "callable": {}, "chr": {}, "dir": {},
		"divmod": {}, "eval": {}, "exec": {}, "filter": {}, "format": {}, "getattr": {},
		"hasattr": {}, "hash": {}, "hex": {}, "id": {}, "iter": {}, "max": {}, "min": {},
		"next": {}, "ord": {}, "pow": {}, "property": {}, "repr": {}, "reversed": {},
		"round": {}, "setattr": {}, "slice": {}, "sorted": {}, "staticmethod": {},
		"classmethod": {}, "sum": {}, "vars": {}, "zip": {},
	},
	"js": {
		"console": {}, "document": {}, "window": {}, "process": {}, "require": {},
		"module": {}, "exports": {}, "toString": {}, "valueOf": {}, "parseInt": {},
		"parseFloat": {}, "encodeURIComponent": {}, "decodeURIComponent": {},
		"Math": {}, "JSON": {}, "Object": {}, "Array": {}, "String": {}, "Number": {},
		"Boolean": {}, "Promise": {}, "Error": {}, "Date": {}, "RegExp": {}, "Map": {},
		"Set": {}, "Symbol": {}, "Proxy": {}, "Reflect": {}, "setTimeout": {},
		"clearTimeout": {}, "setInterval": {}, "clearInterval": {}, "fetch": {},
	},
	"jsx": {
		"console": {}, "document": {}, "window": {}, "process": {}, "require": {},
		"module": {}, "exports": {}, "toString": {}, "valueOf": {}, "parseInt": {},
		"parseFloat": {}, "encodeURIComponent": {}, "decodeURIComponent": {},
		"Math": {}, "JSON": {}, "Object": {}, "Array": {}, "String": {}, "Number": {},
		"Boolean": {}, "Promise": {}, "Error": {}, "Date": {}, "RegExp": {}, "Map": {},
		"Set": {}, "Symbol": {}, "Proxy": {}, "Reflect": {}, "setTimeout": {},
		"clearTimeout": {}, "setInterval": {}, "clearInterval": {}, "fetch": {},
		"React": {}, "useState": {}, "useEffect": {}, "useContext": {}, "useReducer": {},
		"useCallback": {}, "useMemo": {}, "useRef": {},
	},
	"go": {
		"main": {}, "init": {}, "make": {}, "len": {}, "cap": {}, "append": {},
		"copy": {}, "close": {}, "delete": {}, "panic": {}, "recover": {}, "new": {},
		"error": {}, "print": {}, "println": {}, "real": {}, "imag": {}, "complex": {},
	},
	"java": {
		"main": {}, "toString": {}, "equals": {}, "hashCode": {}, "compareTo": {},
		"clone": {}, "finalize": {}, "getClass": {}, "notify": {}, "notifyAll": {},
		"wait": {}, "String": {}, "Object": {}, "Integer": {}, "Double": {},
		"Float": {}, "Boolean": {}, "System": {}, "Math": {},
	},
}

// IsBuiltinSymbol performs an O(1) set lookup to check if a symbol is a language built-in or dunder method.
func IsBuiltinSymbol(langExt, symbol string) bool {
	ext := strings.TrimPrefix(strings.ToLower(langExt), ".")
	if set, ok := BuiltinSymbolSet[ext]; ok {
		_, exists := set[symbol]
		return exists
	}
	return false
}

// IsIgnoredSymbol checks if a symbol is in the user-configured ignored symbols list.
func IsIgnoredSymbol(symbol string, userIgnored []string) bool {
	for _, ign := range userIgnored {
		ign = strings.TrimSpace(ign)
		if ign == "" {
			continue
		}
		if ign == symbol || strings.HasPrefix(symbol, ign) {
			return true
		}
	}
	return false
}
