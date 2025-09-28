<script lang="ts">
  /*
   * Copyright (c) 2025 Bouali Consulting Inc.
   * Author: Kaiss Bouali (kaissb)
   * Company: Bouali Consulting Inc.
   * GitHub: https://github.com/kaissb
   */
  import { onMount, onDestroy } from 'svelte';
  import { Chart, registerables } from 'chart.js';
  Chart.register(...registerables);

  export let containerId: string;
  let canvas: HTMLCanvasElement;
  let chart: Chart;
  let socket: WebSocket;

  onMount(() => {
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    chart = new Chart(ctx, {
      type: 'line',
      data: {
        labels: [],
        datasets: [
          {
            label: 'CPU %',
            data: [],
            borderColor: 'rgb(75, 192, 192)',
            yAxisID: 'y',
          },
          {
            label: 'Memory %',
            data: [],
            borderColor: 'rgb(255, 99, 132)',
            yAxisID: 'y',
          },
        ],
      },
      options: {
        responsive: true,
        scales: {
          y: {
            type: 'linear',
            display: true,
            position: 'left',
            min: 0,
            max: 100,
            ticks: {
              callback: (value) => `${value}%`,
            },
          },
        },
        animation: {
          duration: 0, // General animation time
        },
      },
    });

    socket = new WebSocket(`ws://localhost:8080/ws/stats/${containerId}`);

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      const now = new Date();
      const label = `${now.getHours()}:${now.getMinutes()}:${now.getSeconds()}`;

      chart.data.labels.push(label);
      (chart.data.datasets[0].data as number[]).push(data.cpu_percent);
      (chart.data.datasets[1].data as number[]).push(data.memory_percent);

      // Limit the number of data points
      if (chart.data.labels.length > 30) {
        chart.data.labels.shift();
        (chart.data.datasets[0].data as number[]).shift();
        (chart.data.datasets[1].data as number[]).shift();
      }

      chart.update();
    };

    socket.onerror = (error) => {
      console.error(`WebSocket Error for ${containerId}:`, error);
    };
  });

  onDestroy(() => {
    if (socket) {
      socket.close();
    }
    if (chart) {
      chart.destroy();
    }
  });
</script>

<div class="chart-container">
  <canvas bind:this={canvas}></canvas>
</div>

<style>
  .chart-container {
    position: relative;
    height: 200px;
    width: 100%;
  }
</style>
