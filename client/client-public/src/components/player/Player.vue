<template>
<div class="player-wrapper">
    <div class="player-controlls">
        <NextTrackSVG class="fliped-svg"/>
        <PlayButton/>
        <NextTrackSVG class="player-svg"/>
        <MixButtonSVG class="mix-btn"/>
        <LoopButtonSVG class="loop-btn"/>
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
        <TrackProgressBar/>
    </div>
    <div class="player-volume-wrapper">
        <div class="player-volume-btn"></div>
    </div>
</div>
</template>

<script setup>
import LoopButtonSVG from '../../../public/LoopButtonSVG.vue';
import MixButtonSVG from '../../../public/music/MixButtonSVG.vue';
import NextTrackSVG from '../../../public/NextTrackSVG.vue';
import PlayButton from './PlayButton.vue';
import TrackProgressBar from './TrackProgressBar.vue';

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
    font-size: 14px;
}

.track-artists > span{
    font-size: 12px;
    opacity: .6;
}

.track-image {
    width: 35px;
    height: 35px;
    box-shadow: 1px 2px 4px black;
    overflow: hidden;
    border-radius: 4px;
}

.track-image img{
    width: 35px;
    height: 35px;
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
.mix-btn {
    fill: rgba(255,255,255,.5);
    cursor: pointer;
    transition: all .3s ease-in-out;
    margin-left: 10px;
}
.mix-btn:hover {
    fill: rgba(255,255,255,.8);
    transform: scale(1.1);
}

.loop-btn {
    margin-left: 5px;
    fill: rgba(255,255,255,.5);
    cursor: pointer;
    transition: all .3s ease-in-out;
}

.loop-btn:hover {
    fill: rgba(255,255,255,.8);
    transform: scale(1.1);
}

.player-volume-wrapper {
    width: 100px;
}

.player-volume-svg {

}

.player-volume-btn {
    height: 4px;
    border-radius: 10px;
    background-color: rgba(255,255,255,.3);
}

</style>