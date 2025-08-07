<template>
<div class="user-profile-wrapper">
    <div class="user-profile-main">
        <div class="user-profile-main-photo">
            <div class="user-profile-main-online-status"></div>
            
            
        </div>
        <div class="profile-info">
                <div class="user-profile-name">Артемов Александр</div>
                <div class="user-profile-info-btn"><svg aria-hidden="true" display="block" class="vkuiIcon vkuiIcon--20 vkuiIcon--w-20 vkuiIcon--h-20 vkuiIcon--info_circle_outline_20" width="20" height="20" viewBox="0 0 20 20" fill="none" style="width: 20px; height: 20px;"><path fill="currentColor" fill-rule="evenodd" d="M10 9.25a.75.75 0 0 1 .75.75v4a.75.75 0 0 1-1.5 0v-4a.75.75 0 0 1 .75-.75" clip-rule="evenodd"></path><path fill="currentColor" d="M11 7a1 1 0 1 1-2 0 1 1 0 0 1 2 0"></path><path fill="currentColor" fill-rule="evenodd" d="M3.99 3.99A8.48 8.48 0 0 1 10 1.5c2.347 0 4.473.952 6.01 2.49A8.48 8.48 0 0 1 18.5 10a8.48 8.48 0 0 1-2.49 6.01A8.48 8.48 0 0 1 10 18.5a8.48 8.48 0 0 1-6.01-2.49A8.48 8.48 0 0 1 1.5 10c0-2.347.952-4.473 2.49-6.01M10 3a6.98 6.98 0 0 0-4.95 2.05A6.98 6.98 0 0 0 3 10a6.98 6.98 0 0 0 2.05 4.95A6.98 6.98 0 0 0 10 17a6.98 6.98 0 0 0 4.95-2.05A6.98 6.98 0 0 0 17 10a6.98 6.98 0 0 0-2.05-4.95A6.98 6.98 0 0 0 10 3" clip-rule="evenodd"></path></svg>Подробнее</div>
            </div>
        <div class="user-profile-main-edit"><svg xmlns="http://www.w3.org/2000/svg" width="1024" height="1024" viewBox="0 0 1024 1024"><path fill="#ffffff" d="M257.7 752c2 0 4-.2 6-.5L431.9 722c2-.4 3.9-1.3 5.3-2.8l423.9-423.9a9.96 9.96 0 0 0 0-14.1L694.9 114.9c-1.9-1.9-4.4-2.9-7.1-2.9s-5.2 1-7.1 2.9L256.8 538.8c-1.5 1.5-2.4 3.3-2.8 5.3l-29.5 168.2a33.5 33.5 0 0 0 9.4 29.8c6.6 6.4 14.9 9.9 23.8 9.9m67.4-174.4L687.8 215l73.3 73.3l-362.7 362.6l-88.9 15.7zM880 836H144c-17.7 0-32 14.3-32 32v36c0 4.4 3.6 8 8 8h784c4.4 0 8-3.6 8-8v-36c0-17.7-14.3-32-32-32"/></svg> <span>Редактировать</span></div>
    </div>
    <div class="user-profile-section-mf">
        <div class="section-all"></div>
        <div class="section-friends"></div>
    </div>
</div>
</template>

<script setup>

</script>

<style scoped>

.user-profile-wrapper {
    flex: 1;
    max-width: 1200px;
    height: 100%;
    margin-left: 10px;
}

.user-profile-main {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 15px;
    width: 100%;
    height: 200px;
    border-radius: 8px;
    background-color: #17212b;
    box-shadow: 0 4px 8px black;
}

.user-profile-main-photo {
    position: relative;
    width: 160px;
    height: 160px;
    background-color: rgb(119, 78, 78);
    border-radius: 50%;
    box-shadow: 2px 4px 8px black;
}
.user-profile-main-online-status {
    position: absolute;
    width: 25px;
    right: 10px;
    bottom: 10px;
    height: 25px;
    background-color: rgb(22, 255, 100);
    border-radius: 50%;
    border: 4px solid #17212b;
}

.profile-info {
    display: flex;
    flex-direction: column;
    flex: 1;
    height: 100%;
    justify-content: center;
    margin-left: 15px;
}

.user-profile-name {
    opacity: .9;
    font-size: 22px;
    text-wrap: nowrap;
    font-weight: bold;
}

.user-profile-info-btn {
    position: relative;
    display: flex;
    align-items: center;
    margin-top: 10px;
    opacity: .4;
    font-size: 15px;
    cursor: pointer;
    transition: all .2s ease-in-out;
}

.user-profile-info-btn::after {
    content: "";
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    width: 0;
    margin-left: 24px;
    height: 1px;
    background-color: rgb(255, 255, 255);
    opacity: 0;
    transition: all .2s ease-in-out;
}

.user-profile-info-btn:hover::after {
    width: 80px;
    opacity: .8;
}

.user-profile-info-btn:hover {
    opacity: .8;
}

.user-profile-info-btn > svg {
    margin-right: 4px;
}

.user-profile-main-edit {
    display: flex;
    align-items: center;
    margin-right: 10px;
    padding: 10px 20px;
    cursor: pointer;
    border-radius: 8px;
    background-color: #0e1621;
    transition: all .2s ease-in-out;
    box-shadow: 1px 2px 4px black;
    font-weight: bold;
}

.user-profile-main-edit > svg {
    width: 24px;
    height: 24px;
    fill: rgba(255, 255, 255, 0.6);
    transition: all .2s ease-in-out;
}

.user-profile-main-edit > span {
    color: #ffffff98;
    margin-left: 5px;
    transition: all .2s ease-in-out;
}

.user-profile-main-edit:hover {
    background-color: #ffffff15;
    color: #fffffff1;
}

.user-profile-main-edit:hover > svg {
    fill: #ffffffde;
}

.user-profile-main-edit:hover > span {
    color: rgba(255,255,255, .9);
}

.user-profile-section-mf {
    display: flex;
    width: 100%;
    margin-top: 20px;
}

.section-all {
    flex: 1;
    border-radius: 8px;
    height: 400px;
    background-color: #17212b;
    box-shadow: 2px 4px 8px black;
    margin-right: 20px;
}

.section-friends {
    width: 350px;
    height: 400px;
    background-color: #17212b;
    box-shadow: 2px 4px 8px black;
    border-radius: 8px;
}

</style>