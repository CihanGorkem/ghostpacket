<script lang="ts">
  import { onMount } from 'svelte';
  import Map from './Map.svelte'; // Harita bileşenini çağır

  let attacks: any[] = [];
  let connectionStatus = "Bağlanıyor...";

  onMount(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
      connectionStatus = "🟢 BAĞLANDI: Sistem Online";
    };

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      // Ekranda aynı anda en fazla 15 saldırı göster (Performans için)
      attacks = [...attacks, data].slice(-15); 
      
      // 2 saniye sonra saldırıyı ekrandan sil (Animasyon bitince)
      setTimeout(() => {
        attacks = attacks.filter(a => a !== data);
      }, 2000);
    };

    socket.onclose = () => connectionStatus = "🔴 KOPUK";
  });
</script>

<main>
  <div class="header">
    <h1>GHOST-PACKET <span class="subtitle">Global Cyber Threat Map</span></h1>
    <div class="status">{connectionStatus}</div>
  </div>

  <Map {attacks} />

  <div class="stats">
    <p>Live Attacks: {attacks.length}</p>
  </div>
</main>

<style>
  :global(body) {
    background-color: #050505; /* Simsiyah uzay boşluğu */
    color: #00ff41;
    font-family: 'Courier New', Courier, monospace;
    margin: 0;
    overflow: hidden; /* Kaydırma çubuğunu gizle */
  }
  .header {
    text-align: center;
    padding: 20px;
    z-index: 10;
    position: relative;
  }
  .subtitle {
    font-size: 0.5em;
    color: #666;
    letter-spacing: 2px;
  }
  .stats {
    position: fixed;
    bottom: 20px;
    right: 20px;
    background: rgba(0,0,0,0.8);
    padding: 10px;
    border: 1px solid #333;
  }
</style>