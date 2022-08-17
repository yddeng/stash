package regexp

import (
	"fmt"
	"testing"
)

func TestMatchQQ(t *testing.T) {
	fmt.Println(MatchQQ("1234"))          // 4 f
	fmt.Println(MatchQQ("0123456"))       // 0 f
	fmt.Println(MatchQQ("123456"))        // 6 t
	fmt.Println(MatchQQ("248244153"))     // 9 t
	fmt.Println(MatchQQ("1427896322"))    // 10 t
	fmt.Println(MatchQQ("12345678911"))   // 11 t
	fmt.Println(MatchQQ("123456789111"))  // 12 f
	fmt.Println(MatchQQ("sdfsfdsd"))      // s f
	fmt.Println(MatchQQ("sdfsf12"))       // s f
	fmt.Println(MatchQQ("eee"))           // s f
	fmt.Println(MatchQQ("12348911ee"))    // s f
	fmt.Println(MatchQQ("11234567891ee")) // s f
}

func TestMatchWX(t *testing.T) {
	fmt.Println(MatchWX("1sdfsaf"))               // f
	fmt.Println(MatchWX("ssdf456sf_+"))           // f
	fmt.Println(MatchWX("ssdf456sf_"))            // t
	fmt.Println(MatchWX("A12345678912345678900")) // f
	fmt.Println(MatchWX("A1234567891234567890"))  // t
	fmt.Println(MatchWX("248244154@qq.com"))      // f
}
