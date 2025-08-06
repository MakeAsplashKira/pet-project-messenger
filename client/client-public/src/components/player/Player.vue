<template>
<div class="player-wrapper">
    <div class="player-controlls">
        <NextTrackSVG class="fliped-svg"/>
        <PlayButton/>
        <NextTrackSVG/>
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
             <div class="track-duration">{{ trackData.duration }}</div>
        </div>
        <div class="player-progress-bar-wrapper">
            <div class="player-progress-bar"></div>
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
    gap: 15px;
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

.player-progress-bar-wrapper {
    padding: 3px 0;
    cursor: pointer;
}

.player-progress-bar {
    margin-top: 10px;
    width: 100%;
    height: 4px;
    background-color: rgba(255, 255, 255, 0.356);
    border-radius: 10px;
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
    flex: 1;
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
    width: 40px;
    height: 40px;
    box-shadow: 3px 3px 6px black;
}

.track-image img{
    width: 40px;
    height: 40px;
}

.fliped-svg {
    transform: rotate(180deg);
}

.player-controlls {
    display: flex;
    gap: 5px;
    align-items: center;
}
</style>