<script setup lang="ts">
import Player from './player/Player.vue';

</script>

<template>
    <div class="music-wrapper">
        <section class="player-wrapper">
            <Player class="player"/>
        </section>
    </div>
</template>

<style scoped>
.music-wrapper {
    width: 100%;
}

.player-wrapper {
    background-color: #17212b;
    border-radius: 8px;
    box-shadow: 2px 2px 6px black;
}
.player {
    padding: 10px;
}
</style>