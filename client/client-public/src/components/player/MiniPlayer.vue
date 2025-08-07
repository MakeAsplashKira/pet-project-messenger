<template>
<div class="player-wrapper">
    <div class="player-controlls">
        <NextTrackSVG class="fliped-svg"/>
        <PlayButton/>
        <NextTrackSVG class="player-svg"/>
    </div>
    <div class="player-content">
        <div class="player-main">
            <div class="track-image">
                <img src="../../../public/music/track1.jpg" alt="Track">
            </div>
            <div class="track-info">
                 <div class="track-name">{{trackData.name}}</div>
                 <div class="track-artists">
                    <span v-for="artist in trackData.artists">
                        {{ artist.name }}
                    </span>
                 </div>
             </div>
        </div>
    </div>
</div>
</template>

<script setup>
import NextTrackSVG from '../../../public/NextTrackSVG.vue';
import PlayButton from './PlayButton.vue';

const trackData = 
    {
        name: 'Раскаянье',
        duration: '2:36',
        artists: [
            {
                id: 1,
                name: 'Fovella',
            },
            {
                id: 2,
                name: 'Артём Молоко',
            }   
        ],
        src: ''
    }

</script>

<style scoped>

.player-wrapper {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
}

.player-content {
    display: flex;
    flex-direction: column;
}

.player-main {
    position: relative;
    display: flex;
    align-items: center;
    width: 350px;
}



.track-duration {
    font-size: 14px;
    opacity: .8;
    user-select: none;
    pointer-events: none;
}

.track-info {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    margin-left: 15px;
}

.track-name {
    opacity: .9;
    font-size: 16px;
}

.track-artists > span{
    font-size: 14px;
    opacity: .6;
}

.track-image {
    width: 50px;
    height: 50px;
    box-shadow: 3px 3px 6px black;
}

.track-image img{
    width: 50px;
    height: 50px;
}

.player-controlls {
    display: flex;
    gap: 5px;
    align-items: center;
}

.fliped-svg {
    transform: rotate(180deg);
}
.fliped-svg:hover {
    transform: rotate(180deg) scale(1.05) translateY(-1px);
}

</style>