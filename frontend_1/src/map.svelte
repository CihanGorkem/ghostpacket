<script lang="ts">
  import { onMount } from "svelte";
  import * as d3 from "d3-geo";
  import * as topojson from "topojson-client";

  export let attacks: any[] = [];

  let worldData: any = [];
  
  // Harita projeksiyon ayarı
  const projection = d3.geoNaturalEarth1()
    .scale(160)
    .translate([400, 250]); 

  const pathGenerator = d3.geoPath().projection(projection);

  onMount(async () => {
    const response = await fetch('https://cdn.jsdelivr.net/npm/world-atlas@2/countries-110m.json');
    const data = await response.json();
    worldData = topojson.feature(data, data.objects.countries).features;
  });
</script>

<div class="map-container">
  <svg viewBox="0 0 800 500">
    <defs>
      <filter id="glow">
        <feGaussianBlur stdDeviation="2.5" result="coloredBlur"/>
        <feMerge>
          <feMergeNode in="coloredBlur"/>
          <feMergeNode in="SourceGraphic"/>
        </feMerge>
      </filter>

      <linearGradient id="grad-red" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#ff0055" stop-opacity="0" />
        <stop offset="50%" stop-color="#ff0055" stop-opacity="1" />
        <stop offset="100%" stop-color="#ff0055" stop-opacity="0" />
      </linearGradient>
      <linearGradient id="grad-yellow" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#ffcc00" stop-opacity="0" />
        <stop offset="50%" stop-color="#ffcc00" stop-opacity="1" />
        <stop offset="100%" stop-color="#ffcc00" stop-opacity="0" />
      </linearGradient>
    </defs>

    <g class="countries" filter="url(#glow)">
      {#each worldData as country}
        <path 
          d={pathGenerator(country)}
          class="country-path" 
          role="img"
        />
      {/each}
    </g>

    <g class="attacks">
      {#each attacks as attack}
        {#if attack.src_lon && attack.dst_lon}
          {@const source = projection([attack.src_lon, attack.src_lat])}
          {@const target = projection([attack.dst_lon, attack.dst_lat])}
          
          {#if source && target}
            {@const gradId = (attack.type === 'MALWARE' || attack.type === 'SQL_INJECTION') ? 'url(#grad-yellow)' : 'url(#grad-red)'}
            <path 
              d={`M${source[0]},${source[1]} Q ${(source[0]+target[0])/2}, ${(source[1]+target[1])/2 - 100} ${target[0]},${target[1]}`}
              class="attack-line"
              stroke={gradId}
            />
            
            <g transform={`translate(${target[0]}, ${target[1]})`}>
              <circle r="2" class="impact-core {attack.type}" />
              <circle r="10" class="impact-wave {attack.type}" style="animation-delay: 0s" />
              <circle r="10" class="impact-wave {attack.type}" style="animation-delay: 0.3s" />
            </g>
          {/if}
        {/if}
      {/each}
    </g>
  </svg>
</div>


<style>
  .map-container {
    width: 100%;
    max-width: 1100px;
    margin: 0 auto;
    /* Tüm haritaya genel bir yeşil siber parlama */
    filter: drop-shadow(0 0 15px rgba(0, 255, 65, 0.2)); 
  }
  
  /* Ülkelerin Stili */
  .country-path {
    fill: #0d1117;       /* Daha koyu zemin */
    stroke: #21262d;     /* Çok silik sınırlar */
    stroke-width: 0.5;
    transition: all 0.3s ease;
    cursor: crosshair;
  }

  /* HOVER EFEKTİ: Mouse üzerine gelince ülke parlasın */
  .country-path:hover {
    fill: #1c2128;
    stroke: #00ff41; /* Neon yeşil sınır */
    stroke-width: 1;
    filter: drop-shadow(0 0 5px #00ff41); /* Ülkeye özel parlama */
  }

  /* GÜNCELLENMİŞ SALDIRI ÇİZGİSİ STİLİ */
  /* GÜNCELLENMİŞ ZARİF STİL */
  .attack-line {
    fill: none;
    stroke-width: 1.5; /* Çok daha ince ve zarif */
    stroke-linecap: round;
    opacity: 0.8; /* Hafif şeffaf */
    
    /* 50: Işığın uzunluğu (kısa bir mermi gibi)
       400: Boşluk uzunluğu (uzun bir boşluk)
       Bu sayede ekranda çizgi karmaşası olmaz, tek tek mermiler akar.
    */
    stroke-dasharray: 50 400; 
    
    /* Hızı 2 saniyeye çektik, daha izlenebilir olsun */
    animation: flow 2s linear infinite;
  }
  
  /* Animasyonun bitiş noktasını yeni uzunluğa (50+400=450) göre ayarla */
  @keyframes flow {
    to { 
        stroke-dashoffset: -450; 
    }
  }

  /* Patlama Efektleri */
  .impact-core {
    fill: #fff;
    filter: drop-shadow(0 0 5px #fff); /* Çekirdek parlaması */
  }
  .impact-wave {
    fill: none;
    stroke-width: 2;
    opacity: 0;
    animation: wavePulse 1.5s ease-out infinite;
  }

  /* Renk Atamaları */
  .impact-core.DDOS, .impact-core.BRUTE_FORCE { fill: #ff0055; }
  .impact-wave.DDOS, .impact-wave.BRUTE_FORCE { stroke: #ff0055; }

  .impact-core.MALWARE, .impact-core.SQL_INJECTION { fill: #ffcc00; }
  .impact-wave.MALWARE, .impact-wave.SQL_INJECTION { stroke: #ffcc00; }

  /* GÜNCELLENMİŞ ANİMASYONLAR */
  /* Sürekli akış animasyonu. 
    'to' değeri, stroke-dasharray toplamının negatifi olmalı (-280).
    Bu, çizginin sürekli ileri aktığı yanılsamasını yaratır.
  */
  @keyframes flow {
    to { 
        stroke-dashoffset: -280; 
    }
  }
  
  /* Radar dalgası animasyonu - Biraz daha belirgin */
  @keyframes wavePulse {
    0% { r: 2; opacity: 0.9; stroke-width: 2; }
    100% { r: 35; opacity: 0; stroke-width: 0.5; }
  }
</style>