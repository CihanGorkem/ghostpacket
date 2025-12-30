package main

import (
	
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// 1. Veri Modeli: Bir saldırı paketinin içinde ne olacak?
type AttackData struct {
	SourceLat float64 `json:"src_lat"` // Saldıran Enlem
	SourceLon float64 `json:"src_lon"` // Saldıran Boylam
	DestLat   float64 `json:"dst_lat"` // Hedef Enlem
	DestLon   float64 `json:"dst_lon"` // Hedef Boylam
	Type      string  `json:"type"`    // Saldırı Tipi (DDoS, SQLi, vb.)
}

// WebSocket ayarları (CORS izni veriyoruz ki frontend bağlanabilsin)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// HTTP bağlantısını WebSocket'e yükselt
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Upgrade hatası:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Yeni bir ajan (client) bağlandı!")

	// Sonsuz döngü: Sürekli veri gönder
	for {
		// Rastgele koordinat üret (Basit simülasyon)
		// Gerçek dünya koordinat aralıkları: Lat -90/90, Lon -180/180
		attack := AttackData{
			SourceLat: (rand.Float64() * 180) - 90,
			SourceLon: (rand.Float64() * 360) - 180,
			DestLat:   (rand.Float64() * 180) - 90,
			DestLon:   (rand.Float64() * 360) - 180,
			Type:      getRandomAttackType(),
		}

		// Veriyi JSON'a çevir ve gönder
		err := conn.WriteJSON(attack)
		if err != nil {
			log.Println("Yazma hatası:", err)
			break
		}

		// Hızı ayarla: Her 100 milisaniyede bir saldırı (Makineli tüfek gibi)
		time.Sleep(100 * time.Millisecond)
	}
}

func getRandomAttackType() string {
	types := []string{"DDOS", "SQL_INJECTION", "MALWARE", "BRUTE_FORCE"}
	return types[rand.Intn(len(types))]
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println("Ghost-Packet Sunucusu 8080 portunda dinliyor... (Saldırı simülasyonu hazır)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}