// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"time"
)

// item 구조체 정의
type item struct {
	name   string
	price  int
	amount int
}

// buyer 구조체 정의
type buyer struct {
	point          int
	shoppingBucket map[string]int
}

// 생성자 함수
func newBuyer() *buyer { // 포인터 구조체를 반환함
	d := buyer{} // 구조체 객체를 생성하고 초기화
	d.point = 1000000
	d.shoppingBucket = map[string]int{} // shoppingBucket 필드의 맵을 초기화
	return &d                           // 초기화 한 포인터 구조체를 반환함

}

// item 구조체, buyer 구조체 사용
func buying(itm []item, byr *buyer, itmchoice int, dlt []delivery, d chan bool, numbuy *int, temp map[string]int) {
	// numbuy 를 이용하여 배송 한도를 관리한다

	// panic(), recover(), defer() - 프로그램이 바로 종료되지 않게 하기 위함
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r, "\n")
		}
	}()

	inputamount := 0 // 구매 수량

	fmt.Print("수량을 입력하시오 :")
	fmt.Scanln(&inputamount)
	fmt.Println()

	if inputamount <= 0 {
		panic("올바른 수량을 입력하세요.")
	}
	if byr.point < itm[itmchoice-1].price*inputamount || itm[itmchoice-1].amount < inputamount { // 수량, 포인트로 구매 가능 여부
		panic("주문이 불가능합니다")
	} else {
		for {
			buy := 0 // 살지 장바구니에 담을지 선택
			fmt.Println("1. 바로 주문\n2. 장바구니에 담기")
			fmt.Print("실행할 기능으 입력하시오 :")
			fmt.Scanln(&buy)
			fmt.Println()

			if buy == 1 { // 바로 주문
				if *numbuy < 5 {
					itm[itmchoice-1].amount -= inputamount
					byr.point -= itm[itmchoice-1].price * inputamount
					temp[itm[itmchoice-1].name] = inputamount // 임시저장

					d <- true

					*numbuy++

					fmt.Println("상품이 주문 접수 되었습니다.")

					break

				} else {
					fmt.Println("배송 한도를 초과했습니다. 배송이 완료되면 주문하세요.")
					break
				}
			} else if buy == 2 { // 장바구니에 담기
				checkbucket := false // 중복 물품을 체크하기 위한 변수

				for itms := range byr.shoppingBucket { // 물품 체크
					if itms == itm[itmchoice-1].name {
						checkbucket = true
					}
				}

				if checkbucket == true { // 장바구니에 중복되는 물품이 있을 때
					// 현재 추가하려는 품목 수량 + 기존 장바구니에 담긴 수량 > 잔여 수량 예외 처리
					temp := byr.shoppingBucket[itm[itmchoice-1].name] + inputamount // 수량만 더함
					if temp > itm[itmchoice-1].amount {
						fmt.Println("물품의 잔여 수량을 초과했습니다.")
						break
					}
				} else { // 장바구니에 중복되는 물품이 없을 때
					byr.shoppingBucket[itm[itmchoice-1].name] = inputamount // 새로운 품목 추가
				}

				fmt.Println("상품이 장바구니에 추가되었습니다.")
				break // 구매 for문을 빠져나감
			} else {
				fmt.Println("잘못된 입력입니다. 다시 입력해주세요.")
			}
		}
	}
}

func emptyBucket(byr *buyer) {
	// 장바구니 비었는지 확인, 안비었으면 물품 출력

	if len(byr.shoppingBucket) == 0 {
		fmt.Println("장바구니가 비었습니다.")
	} else {
		for index, val := range byr.shoppingBucket {
			fmt.Printf("%s, 수량: %d\n", index, val)
		}
	}
	fmt.Println()
}

func requiredPoint(itm []item, byr *buyer) (canbuy bool) {
	// 1. 장바구니 총 필요 마일리지 계산
	// 2. 필요 마일리지와 보유 마일리지 출력
	// 3. 마일리지가 부족할 시에는 메세지와 false 반환

	bucketpoint := 0
	for index, val := range byr.shoppingBucket { // 총 필요 마일리지 계산
		for i := 0; i < len(itm); i++ {
			if itm[i].name == index {
				bucketpoint += itm[i].price * val
			}
		}
	}

	fmt.Println("필요 마일리지 :", bucketpoint)
	fmt.Println("보유 마일리지 :", byr.point)
	fmt.Println()
	if byr.point < bucketpoint {
		fmt.Printf("마일리지가 %d점 부족합니다.", bucketpoint-byr.point)
		return false
	}
	return true
}

// 장바구니의 물품 수량이 잔여 수량을 초과할 경우를 판별하는 함수
func excessAmount(itm []item, byr *buyer) (canbuy bool) {

	for index, val := range byr.shoppingBucket {
		for i := 0; i < len(itm); i++ {
			if itm[i].name == index { // itm[i] 에 물품명 입력
				if itm[i].amount < val { // 장바구니의 물품 총 개수가 판매하는 물품 개수보다 클때
					fmt.Printf("%s, %d개 초과", itm[i].name, val-itm[i].amount)
					return false
				}
			}
		}
	}
	return true
}

// 주문 기능
func bucketBuying(itm []item, byr *buyer, numbuy *int, temp map[string]int, d chan bool) {
	// 품목의 수량을 장바구니 수량만큼 차감하고 수용자 포인트도 차감했다면 장바구니를 초기화하고 주문 접수
	// numbuy 를 이용하여 배송 한도를 관리한다
	defer func() { // 함수 내에서 장바구니의 물품이 존재하지 않는데 주문할 경우 panic으로 발생시키고 복구
		if r := recover(); r != nil {
			fmt.Println("\n", r, "\n")
		}
	}()

	if len(byr.shoppingBucket) == 0 {
		panic("주문 가능한 목록이 없습니다.")
	} else {

		if *numbuy < 5 {
			for index, val := range byr.shoppingBucket {
				temp[index] = val // 임시 저장

				for i := range itm {
					if itm[i].name == index {
						itm[i].amount -= val            // 수량 차감
						byr.point -= itm[i].price * val // 포인트 차감
					}
				}
			}
			d <- true // 배송 시작

			byr.shoppingBucket = map[string]int{} // 장바구니 초기화
			*numbuy++

			fmt.Println("주문 접수 되었습니다.")
		} else {
			fmt.Println("배송 한도를 초과했습니다. 배송이 완료되면 주문하세요.")
		}

		byr.shoppingBucket = map[string]int{} // 장바구니 초기화
		fmt.Println("주문 접수 되었습니다.")
	}
}

// delivery 구조체
type delivery struct {
	// 객체를 트럭 한대라고 생각
	// 5대의 트럭이 필요하면 5개의 객체를 생성

	status      string         // 배송상태
	onedelivery map[string]int // 한번에 배송하는 물품의 뜻으로 생각
}

// 생성자 함수
func newDelivery() delivery {
	d := delivery{}
	d.onedelivery = map[string]int{}
	return d
}

// 배송 상태 확인 고루틴 생성
func deliveryStatus(d chan bool, i int, deliverylist []delivery, num *int, temp *map[string]int) {
	for {
		if <-d {
			for index, val := range *temp {
				deliverylist[i].onedelivery[index] = val // 임시 저장한 데이터를 배송 상품에 저장
			}

			*temp = map[string]int{} // 임시 데이터 초기화

			deliverylist[i].status = "주문접수"
			time.Sleep(time.Second * 10)

			deliverylist[i].status = "배송중"
			time.Sleep(time.Second * 30)

			deliverylist[i].status = "배송완료"
			time.Sleep(time.Second * 10)

			deliverylist[i].status = ""
			*num--
			deliverylist[i].onedelivery = map[string]int{} // 배송 리스트에서 물품 지우기
		}
	}
}

func main() {

	// 주문 건수를 제한하는 변수로 주문 건수 5개로 제한
	// 주문처리가 되면 배송 건수를 1 증가시키고, 5 초과되면 주문하지 못하도록 한다.
	numbuy := 0 // 주문한 개수

	// item 슬라이스 생성
	items := make([]item, 5) // 물품 목록

	// buyer 객체 생성 생성자를 이용하여
	buyer := newBuyer() // 구매자 정보(장바구니, 마일리)

	// delivery 슬라이스 생성 (객체 5개 생성)
	deliverylist := make([]delivery, 5) // 배송 중인 상품 목록

	deliverystart := make(chan bool) // 주문 시작 신호 송수신 채널

	for i := 0; i < 5; i++ { // 배송 상품 객체 5개 생성
		deliverylist[i] = newDelivery()
	}

	// 주문한 품목과 수량을 저장하는 map tempdelivery  생성
	tempdelivery := make(map[string]int) // 배달 물품 임시 저장

	for i := 0; i < 5; i++ {
		// deliverystart -> 채널 송/수신을 위한 것
		// i -> 각각의 고루틴(트럭)을 구분하기 위해 전달
		// deliverylist -> 배송 객체를 받고 마지막으로 numbuy 는 최대 배송 가능 횟수를 포인터 변수로 전달
		// 임시 저장한 주문품목, 수량을 func deliveryStatus() 에 가져와 deliverylist[i].onedelivery에 저장한다. 그리고 값 초기화
		time.Sleep(time.Millisecond) // 고루틴 순서대로 실행되도록 약간의 딜레이
		go deliveryStatus(deliverystart, i, deliverylist, &numbuy, &tempdelivery)
	}

	items[0] = item{"텀블러", 10000, 30}
	items[1] = item{"롱패딩", 500000, 20}
	items[2] = item{"백패", 400000, 20}
	items[3] = item{"운동화", 150000, 50}
	items[4] = item{"빼빼로", 1200, 500}

	// 프로그램을 종료하기 전까지 계속 반복
	for {
		menu := 0 // 첫메뉴

		fmt.Println("1. 구매")
		fmt.Println("2. 잔여 수량 확인")
		fmt.Println("3. 잔여 마일리지 확인")
		fmt.Println("4. 배송 상태 확인")
		fmt.Println("5. 장바구니 확인")
		fmt.Println("6. 프로그램 종료")
		fmt.Print("실행할 기능을 입력하시오 :")

		fmt.Scanln(&menu)
		fmt.Println()

		if menu == 1 { // 물건 구매
			for {
				itemchoice := 0
				// 물건 리스트 출력
				for i := 0; i < 5; i++ {
					fmt.Printf("물품%d: %s, 가격: %d, 잔여 수량: %d\n", i+1, items[i].name, items[i].price, items[i].amount)
				}

				fmt.Print("구매할 물품을 선택하세요 :")
				fmt.Scanln(&itemchoice)
				fmt.Println()

				// buying() 함수를 이용하여 구매
				if itemchoice == 1 {
					buying(items, buyer, itemchoice, deliverylist, deliverystart, &numbuy, tempdelivery)
					break
				} else if itemchoice == 2 {
					buying(items, buyer, itemchoice, deliverylist, deliverystart, &numbuy, tempdelivery)
					break
				} else if itemchoice == 3 {
					buying(items, buyer, itemchoice, deliverylist, deliverystart, &numbuy, tempdelivery)
					break
				} else if itemchoice == 4 {
					buying(items, buyer, itemchoice, deliverylist, deliverystart, &numbuy, tempdelivery)
					break
				} else if itemchoice == 5 {
					buying(items, buyer, itemchoice, deliverylist, deliverystart, &numbuy, tempdelivery)
					break
				} else {
					fmt.Printf("잘못된 입력입니다. 다시 입력해주세요.\n")
				}
			}

			fmt.Print("엔터를 입력하면 메뉴 화면으로 돌아갑니다")
			fmt.Scanln()
		} else if menu == 2 { // 남은 수량 확인

			// 슬라이스 사용
			for i := 0; i < 5; i++ {
				fmt.Printf("%s, 잔여 수량: %d\n", items[i].name, items[i].amount)
			}
			fmt.Print("엔터를 입력하면 메뉴 화면으로 돌아갑니다.")
			fmt.Scanln()
		} else if menu == 3 { // 잔여 마일리지 확인
			// 생성자함수 사용
			fmt.Printf("현재 잔여 마일리지는 %d점 입니다.\n", buyer.point)
			fmt.Print("엔터를 입력하면 메뉴 화면으로 돌아갑니다.")
			fmt.Scanln()
		} else if menu == 4 { // 배송 상태 확인
			total := 0
			for i := 0; i < 5; i++ {
				total += len(deliverylist[i].onedelivery)
			}
			if total == 0 {
				fmt.Println("배송중인 상품이 없습니다.")
			} else {
				// deliveryStatus 루틴에서  deliverylist[i].status 필드값만 출력하면 된다.
				for i := 0; i < len(deliverylist); i++ {
					if len(deliverylist[i].onedelivery) != 0 { // 배송중인 항목만 출력
						for index, val := range deliverylist[i].onedelivery {
							fmt.Printf("%s %d개 / ", index, val)
						}
						fmt.Printf("배송상황: %s\n", deliverylist[i].status)
					}
				}
			}

			fmt.Print("엔터를 입력하면 메뉴 화면으로 돌아갑니다.")
			fmt.Scanln()
		} else if menu == 5 { // 장바구니 확인
			bucketmenu := 0

			for {
				emptyBucket(buyer) // 장바구니 안이 비었는지 확인하는 함수. 안비었으면 물품 출력

				// 둘중 한개라도 false가 나오면 어차피 구매가 불가능 하기 때문에 똑같은 변수에 반환값을 초기화
				canbuy := requiredPoint(items, buyer) // 살수 있는지 없는지 확인하는 canbuy 선언 및 초기화
				canbuy = excessAmount(items, buyer)   // 장바구니 비었는지 확인, 안비었으면 물품 출력

				fmt.Println("1. 장바구니 상품 주문")
				fmt.Println("2. 장바구니 초기화")
				fmt.Println("3. 메뉴로 돌아가기")

				fmt.Println("실행할 기능을 입력하세요: ")
				fmt.Scanln(&bucketmenu)
				fmt.Println()

				if bucketmenu == 1 {
					if canbuy {
						bucketBuying(items, buyer, &numbuy, tempdelivery, deliverystart)
						break
					} else {
						fmt.Println("구매할 수 없습니다.")
						break
					}
					// 장바구니 상품 주문
				} else if bucketmenu == 2 {
					buyer.shoppingBucket = map[string]int{} // 장바구니 초기화
					fmt.Println("장바구니를 초기화했습니다.")
					break
				} else if bucketmenu == 3 {
					fmt.Println()
					break
				} else {
					fmt.Println("잘못된 입력입니다. 다시 입력해주세요.")
				}
			}

			fmt.Print("엔터를 입력하면 메뉴 화면으로 돌아갑니다.")
			fmt.Scanln()

		} else if menu == 6 { // 프로그램 종료
			fmt.Println("프로그램을 종료합니다.")

			return // main함수 종료
		} else {
			fmt.Printf("잘못된 입력입니다. 다시 입력해주세요.\n")
		}
	}
}
