<template>
  <div 
    class="auth-modal"
    @click.self="closeModal"
  >
    <div class="auth-modal__content">
      <!-- Шаг 1: Email -->
      <div v-if="regStep === 1" class="auth-step">
        <div class="auth-header-1">
            Давай зарегестрируем вас!
        </div>
        <div class="email-verif-code" :class="{'email-verif-code-active': emailApproved && !emailAuthDone}">
              <div class="input-icon"><svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32"><circle cx="25" cy="20" r="1" fill="#ffffff"/><path fill="#ffffff" d="M19.414 30H15v-4.414l5.034-5.034A4.607 4.607 0 0 1 20 20a5 5 0 1 1 4.448 4.966zM17 28h1.586l5.206-5.206l.54.124a3.035 3.035 0 1 0-2.25-2.25l.124.54L17 26.414zM6 8h2v8H6zM2 8h2v8H2zm16 0h2v6h-2zm-4 8h-2a2 2 0 0 1-2-2v-4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v4a2 2 0 0 1-2 2m-2-2h2v-4h-2zM2 18h2v8H2zm12 0h2v4h-2zm-4 8H8a2 2 0 0 1-2-2v-4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v4a2 2 0 0 1-2 2m-2-2h2v-4H8zM2 2h2v4H2zm12 0h2v4h-2zm4 0h2v4h-2zm-8 4H8a2 2 0 0 1-2-2V2h2v2h2V2h2v2a2 2 0 0 1-2 2"/></svg></div>
              <input v-model="code" type="text" class="verify-input" placeholder="code">
              <div class="send-code-container">
                <div @click="sendVerifCode()" class="verif-input-icon"><span v-if="timerCount!= 0">{{ timerCount < 10? '0'+timerCount: timerCount }}</span><svg v-if="timerCount == 0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 48 48"><g :style="{opacity: verifCodeSent? .8 : 0}" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="4"><path d="M44 24V9H4v30h20m7-3l5 4l8-10"/><path d="m4 9l20 15L44 9"/></g><g :style="{opacity: !verifCodeSent? .8 : 0}" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="4"><path d="M44 24V9H4v30h20m20-5H30m9-5l5 5l-5 5"/><path d="m4 9l20 15L44 9"/></g></svg></div>
                <div v-if="verifCodeStatus==-1" class="verify-loading"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#ffffff" d="M12 2A10 10 0 1 0 22 12A10 10 0 0 0 12 2Zm0 18a8 8 0 1 1 8-8A8 8 0 0 1 12 20Z" opacity=".5"/><path fill="#ffffff" d="M20 12h2A10 10 0 0 0 12 2V4A8 8 0 0 1 20 12Z"><animateTransform attributeName="transform" dur="1s" from="0 12 12" repeatCount="indefinite" to="360 12 12" type="rotate"/></path></svg></div>
                <div v-if="verifCodeStatus==1" class="verify-loading"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"><path fill="#35fb32" d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0m-3.97-3.03a.75.75 0 0 0-1.08.022L7.477 9.417L5.384 7.323a.75.75 0 0 0-1.06 1.06L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-.01-1.05z"/></svg></div>
                <div v-if="verifCodeStatus==-2" class="verify-loading"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 42 42"><path fill="#cc2424" d="m29.582 8.683l-.129.12L8.3 29.954a3.308 3.308 0 0 0-.547.688c-2.04-2.639-3.233-6-3.233-9.701c0-8.797 6.626-15.482 15.421-15.482c3.691 0 7.014 1.185 9.641 3.224M10.937 33.704c.189-.117.388-.287.606-.507l21.151-21.151l.041-.04c1.74 2.518 2.746 5.602 2.746 8.994c0 8.785-6.696 15.541-15.481 15.541c-3.432 0-6.546-1.035-9.063-2.837M.5 21C.5 31.775 9.235 40.5 20 40.5c10.767 0 19.501-8.725 19.501-19.5s-8.734-19.5-19.5-19.5S.5 10.225.5 21"/></svg></div>
                <div v-if="verifCodeStatus==0" @click="verifyCode()"  class="email-sent-btn"><svg  aria-hidden="true" display="block" viewBox="0 0 12 8" fill="none"><path fill="currentColor" fill-rule="evenodd" d="M2.156 2.295a.75.75 0 0 1 1.051-.137L6 4.306l2.793-2.148a.75.75 0 1 1 .914 1.189l-3.25 2.5a.75.75 0 0 1-.914 0l-3.25-2.5a.75.75 0 0 1-.137-1.052" clip-rule="evenodd"></path></svg></div>         
              </div>
            </div>
        <div class="auth-content-1">
          <div class="input-container" :class="{'email-auth-end':emailAuthDone }">
            <div class="input-icon"><svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32"><path fill="#ffffff" d="M28 6H4a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h24a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2m-2.2 2L16 14.78L6.2 8ZM4 24V8.91l11.43 7.91a1 1 0 0 0 1.14 0L28 8.91V24Z"/></svg></div>
            <input
            @focus="isFocused = true"
            @blur="isFocused = false"
            type="email"
            formnovalidate
            name="email"
            placeholder="Введите ваш E-mail..."
            v-model="email"
            @input="onEmailInput">
            <div class="email-menu">
              <div v-if="emailStatus===0" class="email-status-icon" :style="{opacity: email.length > 0 && isFocused? 1: 0}"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#ffffff" d="M12 2A10 10 0 1 0 22 12A10 10 0 0 0 12 2Zm0 18a8 8 0 1 1 8-8A8 8 0 0 1 12 20Z" opacity=".5"/><path fill="#ffffff" d="M20 12h2A10 10 0 0 0 12 2V4A8 8 0 0 1 20 12Z"><animateTransform attributeName="transform" dur="1s" from="0 12 12" repeatCount="indefinite" to="360 12 12" type="rotate"/></path></svg></div>   
              <div v-if="emailStatus===1" class="email-status-icon suc-e"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"><path fill="#35fb32" d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0m-3.97-3.03a.75.75 0 0 0-1.08.022L7.477 9.417L5.384 7.323a.75.75 0 0 0-1.06 1.06L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-.01-1.05z"/></svg></div>   
              <div v-if="emailStatus===-1" class="email-status-icon" :style="{opacity: email.length > 0}"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 42 42"><path fill="#cc2424" d="m29.582 8.683l-.129.12L8.3 29.954a3.308 3.308 0 0 0-.547.688c-2.04-2.639-3.233-6-3.233-9.701c0-8.797 6.626-15.482 15.421-15.482c3.691 0 7.014 1.185 9.641 3.224M10.937 33.704c.189-.117.388-.287.606-.507l21.151-21.151l.041-.04c1.74 2.518 2.746 5.602 2.746 8.994c0 8.785-6.696 15.541-15.481 15.541c-3.432 0-6.546-1.035-9.063-2.837M.5 21C.5 31.775 9.235 40.5 20 40.5c10.767 0 19.501-8.725 19.501-19.5s-8.734-19.5-19.5-19.5S.5 10.225.5 21"/></svg></div>   
              <div v-show="!emailAuthDone" @click="sendEmailToVerif()" :style="{opacity: emailReadyToSend? 0.8 : 0.2, pointerEvents: emailReadyToSend? 'all':'none'}" class="email-sent-btn"><svg  aria-hidden="true" display="block" viewBox="0 0 12 8" fill="none"><path fill="currentColor" fill-rule="evenodd" d="M2.156 2.295a.75.75 0 0 1 1.051-.137L6 4.306l2.793-2.148a.75.75 0 1 1 .914 1.189l-3.25 2.5a.75.75 0 0 1-.914 0l-3.25-2.5a.75.75 0 0 1-.137-1.052" clip-rule="evenodd"></path></svg></div>         
            </div>
          </div>
          <div class="input-container password" :style="{marginTop: emailApproved && !emailAuthDone? '76px':'20px'}">
            <div class="input-icon"><svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32"><path fill="#ffffff" d="M21 2a8.998 8.998 0 0 0-8.612 11.612L2 24v6h6l10.388-10.388A9 9 0 1 0 21 2m0 16a7.013 7.013 0 0 1-2.032-.302l-1.147-.348l-.847.847l-3.181 3.181L12.414 20L11 21.414l1.379 1.379l-1.586 1.586L9.414 23L8 24.414l1.379 1.379L7.172 28H4v-3.172l9.802-9.802l.848-.847l-.348-1.147A7 7 0 1 1 21 18"/><circle cx="22" cy="10" r="2" fill="#ffffff"/></svg></div>
            <input placeholder="Введите пароль..."  :type="showPassword? 'text':'password'" name="password" v-model="password">
            <div class="password-show" @click="tooglePassowrdShowing">
              <div class="password-show-icon" :style="{opacity: showPassword? .8 : 0}"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path d="M21.257 10.962c.474.62.474 1.457 0 2.076C19.764 14.987 16.182 19 12 19c-4.182 0-7.764-4.013-9.257-5.962a1.692 1.692 0 0 1 0-2.076C4.236 9.013 7.818 5 12 5c4.182 0 7.764 4.013 9.257 5.962"/><circle cx="12" cy="12" r="3"/></g></svg></div>
              <div class="password-hide-icon" :style="{opacity: !showPassword? .8 : 0}"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2 10s3.5 4 10 4s10-4 10-4M4 11.645L2 14m20 0l-1.996-2.352M8.914 13.68L8 16.5m7.063-2.812L16 16.5"/></svg></div>
            </div>
            <div class="input-status">
              <svg :style="{opacity: password.length>=8 && passwordStrength >= 4 && !isTyping? .9 : 0 }" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"><path fill="#35fb32" d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0m-3.97-3.03a.75.75 0 0 0-1.08.022L7.477 9.417L5.384 7.323a.75.75 0 0 0-1.06 1.06L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-.01-1.05z"/></svg>
              <svg :style="{opacity: (password.length<8 || passwordStrength <4) && password.length>=1 && !isTyping? .9 : 0}"  xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 42 42"><path fill="#cc2424" d="m29.582 8.683l-.129.12L8.3 29.954a3.308 3.308 0 0 0-.547.688c-2.04-2.639-3.233-6-3.233-9.701c0-8.797 6.626-15.482 15.421-15.482c3.691 0 7.014 1.185 9.641 3.224M10.937 33.704c.189-.117.388-.287.606-.507l21.151-21.151l.041-.04c1.74 2.518 2.746 5.602 2.746 8.994c0 8.785-6.696 15.541-15.481 15.541c-3.432 0-6.546-1.035-9.063-2.837M.5 21C.5 31.775 9.235 40.5 20 40.5c10.767 0 19.501-8.725 19.501-19.5s-8.734-19.5-19.5-19.5S.5 10.225.5 21"/></svg>
              <svg :style="{opacity: isTyping? .9:0}"  xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#ffffff" d="M12 2A10 10 0 1 0 22 12A10 10 0 0 0 12 2Zm0 18a8 8 0 1 1 8-8A8 8 0 0 1 12 20Z" opacity=".5"/><path fill="#ffffff" d="M20 12h2A10 10 0 0 0 12 2V4A8 8 0 0 1 20 12Z"><animateTransform attributeName="transform" dur="1s" from="0 12 12" repeatCount="indefinite" to="360 12 12" type="rotate"/></path></svg>
            </div>
            <div class="password-strength-container">
              <div :class="{'bad-11': passwordStrength>=1}" class="password-strength-node bad-1"></div>
              <div :class="{'bad-22': passwordStrength>=2}" class="password-strength-node bad-2"></div>
              <div :class="{'warning-11': passwordStrength>=3}" class="password-strength-node warning-1"></div>
              <div :class="{'warning-22': passwordStrength>=4}" class="password-strength-node warning-2"></div>
              <div :class="{'good-11': passwordStrength>=5}" class="password-strength-node good-1"></div>
              <div :class="{'good-22': passwordStrength>=6}" class="password-strength-node good-2"></div>
              <div :class="{'awesome-11': passwordStrength>=7}" class="password-strength-node awesome-1"></div>
              <div :class="{'awesome-22': passwordStrength>=8}" class="password-strength-node awesome-2"></div>
            </div>
          </div>
          <div class="input-container">
            <div class="input-icon"><svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32"><path style="transform: translate(12px, 17px) scale(0.5);" fill="#ffffff" d="M21 2a8.998 8.998 0 0 0-8.612 11.612L2 24v6h6l10.388-10.388A9 9 0 1 0 21 2m0 16a7.013 7.013 0 0 1-2.032-.302l-1.147-.348l-.847.847l-3.181 3.181L12.414 20L11 21.414l1.379 1.379l-1.586 1.586L9.414 23L8 24.414l1.379 1.379L7.172 28H4v-3.172l9.802-9.802l.848-.847l-.348-1.147A7 7 0 1 1 21 18"/><circle cx="22" cy="10" r="2" fill="#ffffff"/><path fill="#ffffff" d="M21 2a8.998 8.998 0 0 0-8.612 11.612L2 24v6h6l10.388-10.388A9 9 0 1 0 21 2m0 16a7.013 7.013 0 0 1-2.032-.302l-1.147-.348l-.847.847l-3.181 3.181L12.414 20L11 21.414l1.379 1.379l-1.586 1.586L9.414 23L8 24.414l1.379 1.379L7.172 28H4v-3.172l9.802-9.802l.848-.847l-.348-1.147A7 7 0 1 1 21 18"/><circle cx="22" cy="10" r="2" fill="#ffffff"/></svg></div>
            <input v-model="passwordR" :type="showPassword?'text':'password'" placeholder="Введите пароль ещё раз...">
            <div class="password-show" @click="tooglePassowrdShowing">
              <div class="password-show-icon" :style="{opacity: showPassword? .8 : 0}"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path d="M21.257 10.962c.474.62.474 1.457 0 2.076C19.764 14.987 16.182 19 12 19c-4.182 0-7.764-4.013-9.257-5.962a1.692 1.692 0 0 1 0-2.076C4.236 9.013 7.818 5 12 5c4.182 0 7.764 4.013 9.257 5.962"/><circle cx="12" cy="12" r="3"/></g></svg></div>
              <div class="password-hide-icon" :style="{opacity: !showPassword? .8 : 0}"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2 10s3.5 4 10 4s10-4 10-4M4 11.645L2 14m20 0l-1.996-2.352M8.914 13.68L8 16.5m7.063-2.812L16 16.5"/></svg></div>
            </div>
            <div class="input-status">
              <svg :style="{opacity: password==passwordR && password.length>0 && !isTypingR && validatePassword()? .9:0}" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"><path fill="#35fb32" d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0m-3.97-3.03a.75.75 0 0 0-1.08.022L7.477 9.417L5.384 7.323a.75.75 0 0 0-1.06 1.06L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-.01-1.05z"/></svg>
              <svg :style="{opacity: password!=passwordR && passwordR.length>0 && !isTypingR? .9:0}"  xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 42 42"><path fill="#cc2424" d="m29.582 8.683l-.129.12L8.3 29.954a3.308 3.308 0 0 0-.547.688c-2.04-2.639-3.233-6-3.233-9.701c0-8.797 6.626-15.482 15.421-15.482c3.691 0 7.014 1.185 9.641 3.224M10.937 33.704c.189-.117.388-.287.606-.507l21.151-21.151l.041-.04c1.74 2.518 2.746 5.602 2.746 8.994c0 8.785-6.696 15.541-15.481 15.541c-3.432 0-6.546-1.035-9.063-2.837M.5 21C.5 31.775 9.235 40.5 20 40.5c10.767 0 19.501-8.725 19.501-19.5s-8.734-19.5-19.5-19.5S.5 10.225.5 21"/></svg>
              <svg :style="{opacity: isTypingR && passwordR.length>0? .9:0}"  xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#ffffff" d="M12 2A10 10 0 1 0 22 12A10 10 0 0 0 12 2Zm0 18a8 8 0 1 1 8-8A8 8 0 0 1 12 20Z" opacity=".5"/><path fill="#ffffff" d="M20 12h2A10 10 0 0 0 12 2V4A8 8 0 0 1 20 12Z"><animateTransform attributeName="transform" dur="1s" from="0 12 12" repeatCount="indefinite" to="360 12 12" type="rotate"/></path></svg>
            </div>         
          </div>
          <div class="input-container">
            <div class="input-icon username"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#ffffff" fill-rule="evenodd" d="M7.75 7.5a4.25 4.25 0 1 1 8.5 0a4.25 4.25 0 0 1-8.5 0M12 4.75a2.75 2.75 0 1 0 0 5.5a2.75 2.75 0 0 0 0-5.5m-4 10A2.25 2.25 0 0 0 5.75 17v1.188c0 .018.013.034.031.037c4.119.672 8.32.672 12.438 0a.037.037 0 0 0 .031-.037V17A2.25 2.25 0 0 0 16 14.75h-.34a.253.253 0 0 0-.079.012l-.865.283a8.751 8.751 0 0 1-5.432 0l-.866-.283a.252.252 0 0 0-.077-.012zM4.25 17A3.75 3.75 0 0 1 8 13.25h.34c.185 0 .369.03.544.086l.866.283a7.251 7.251 0 0 0 4.5 0l.866-.283c.175-.057.359-.086.543-.086H16A3.75 3.75 0 0 1 19.75 17v1.188c0 .754-.546 1.396-1.29 1.517a40.095 40.095 0 0 1-12.92 0a1.537 1.537 0 0 1-1.29-1.517z" clip-rule="evenodd"/></svg></div>
            <input type="text" placeholder="@Ваш_Идентификатор" v-model="username">
            <div class="input-status">
              <svg :style="{opacity: !isTypingUN && isUNvalidByServer && isUNValid? .9:0}" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"><path fill="#35fb32" d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0m-3.97-3.03a.75.75 0 0 0-1.08.022L7.477 9.417L5.384 7.323a.75.75 0 0 0-1.06 1.06L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-.01-1.05z"/></svg>
              <svg :style="{opacity: !isTypingUN && (!isUNValid || !isUNvalidByServer) && username.length>0? .9:0}"  xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 42 42"><path fill="#cc2424" d="m29.582 8.683l-.129.12L8.3 29.954a3.308 3.308 0 0 0-.547.688c-2.04-2.639-3.233-6-3.233-9.701c0-8.797 6.626-15.482 15.421-15.482c3.691 0 7.014 1.185 9.641 3.224M10.937 33.704c.189-.117.388-.287.606-.507l21.151-21.151l.041-.04c1.74 2.518 2.746 5.602 2.746 8.994c0 8.785-6.696 15.541-15.481 15.541c-3.432 0-6.546-1.035-9.063-2.837M.5 21C.5 31.775 9.235 40.5 20 40.5c10.767 0 19.501-8.725 19.501-19.5s-8.734-19.5-19.5-19.5S.5 10.225.5 21"/></svg>
              <svg :style="{opacity: isTypingUN? .9:0}"  xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#ffffff" d="M12 2A10 10 0 1 0 22 12A10 10 0 0 0 12 2Zm0 18a8 8 0 1 1 8-8A8 8 0 0 1 12 20Z" opacity=".5"/><path fill="#ffffff" d="M20 12h2A10 10 0 0 0 12 2V4A8 8 0 0 1 20 12Z"><animateTransform attributeName="transform" dur="1s" from="0 12 12" repeatCount="indefinite" to="360 12 12" type="rotate"/></path></svg>
            </div>         
          </div>
          <div class="bottom-side">
            <div class="skip-steps-button" :class="{'disabled-btm-btn' : !bottomButtonActive}" @click="completeRegistration">Законить регистрацию</div>
            <div class="next-step-button" :class="{'disabled-btm-btn' : !bottomButtonActive}" @click="goToNextStep">Следующий шаг</div>
          </div>
        </div>
        </div>

      <!-- Шаг 2: Код подтверждения -->
      <div v-if="regStep === 2" class="auth-step">
        <h3>Введите код из письма</h3>
      </div>

      <!-- Шаг 3: Профиль -->
      <div v-if="regStep === 3" class="auth-step">
        <h3>Создайте профиль</h3>
        
      </div>
    </div>
    <div class="auth-modal__steps">
        <div class="auth-steps-left">
   
            <div class="auth-step-anim" v-for="step in registrationSteps">
              <LiquidProgress :text="step.text" :fillPercent="step.fillPercent" :cSize="step.cSize" :current-stage="step.currentStage"/>
            </div>
        </div>
        <div class="auth-anim" ref="lottieContainer"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed, onUnmounted, watchEffect} from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useNotifications } from '@/stores/notifications';
import lottie from 'lottie-web';
import LiquidProgress from './anim/LiquidProgress.vue';
import useNotify from '@/composable/useNotify';
import { checkEmailOnServer, checkUsername, registerUser, sendVerificationCodeOnServer, VerifyCodeOnServer } from '@/api.js';



const {notify} = useNotify()

const lottieContainer = ref(null)
const regStep = ref(1)

onMounted(async () => {
  const animationData = await import('@/assets/auth-anim.json');
  lottie.loadAnimation({
    container: lottieContainer.value,
    renderer: 'svg',
    animationData: animationData.default,
    loop: true,
  });
});

const registrationSteps = ref([
  {text:'1',fillPercent: 0, cSize: 70, currentStage: true},
  {text:'2',fillPercent: 0, cSize: 80, currentStage: false},
  {text:'3',fillPercent: 0, cSize: 90, currentStage: false},
  {text:'Готово!',fillPercent: 0, cSize: 100, currentStage: false},  
])

const goToNextStep=()=> {
  if(bottomButtonActive) {
    updateCurrentStepFillPercentToMax()
    regStep.value+=1
    registrationSteps.value[regStep.value-1].currentStage = true
  }
}

const updateCurrentStepFillPercent = (value) => {
  if (value < 0 && registrationSteps.value[regStep.value-1].fillPercent == 0) {
    registrationSteps.value[regStep.value-1].fillPercent = 0
    return
  }
  registrationSteps.value[regStep.value-1].fillPercent += value
}
const updateCurrentStepFillPercentToMax = () => {
  registrationSteps.value[regStep.value-1].fillPercent = 110
}

const password = ref('')
const passwordR = ref('')

const isTyping = ref(false)
const isTypingR = ref(false)

let timeoutPassword = null
let timeoutPasswordR = null
let passwordAnimUpdated = false
let passwordRAnimUpdated = false

watch(passwordR, ()=> {
  isTypingR.value = true
  clearTimeout(timeoutPasswordR)
  timeoutPasswordR = setTimeout(()=> {
    isTypingR.value = false
    validatePasswordR()
  }, 1000)
})




const validatePasswordR = () => {
    const isLastPasswordValid = validatePassword()
    const isMatch = password.value === passwordR.value
    
    if (isLastPasswordValid && isMatch) {
      if (!passwordRAnimUpdated) {
        updateCurrentStepFillPercent(20)
        passwordRAnimUpdated = true
      }
      return true
    } else {
      if (passwordRAnimUpdated) {
        updateCurrentStepFillPercent(-20)
        passwordRAnimUpdated = false
      }
      return false
    }
}

watch(password, ()=> {
  isTyping.value = true
  clearTimeout(timeoutPassword)
  timeoutPassword = setTimeout(()=> {
    isTyping.value=false
    validatePassword()
    if (passwordR.value.length>0) {
      validatePasswordR()
    }
  }, 1000)
})

const validatePassword = () => {
  const isValid = password.value.length>=8&&passwordStrength.value>=4
    if(isValid && !passwordAnimUpdated) {
        updateCurrentStepFillPercent(20)
        passwordAnimUpdated = true
        return true
    }
    else if (!isValid && passwordAnimUpdated) {
      updateCurrentStepFillPercent(-20)
      passwordAnimUpdated = false
      return false
    }
    return isValid
}

const isUNValid = ref(false)
const isUNvalidByServer = ref(false)
const emailAuthDone = ref(false)
const bottomButtonActive = ref(false)


watchEffect(() => {
  // Явно обращаемся ко всем зависимостям
  const conditions = [
    validatePassword(),
    validatePasswordR(),
    emailAuthDone.value,
    isUNValid.value,
    isUNvalidByServer.value
  ]
  
  bottomButtonActive.value = conditions.every(Boolean)
})
const passwordStrength = computed(() => {
  const pass = password.value;
  if (pass.length === 0) return 0;
  if (/[а-яА-ЯЁё]/.test(pass)) {
    return 0; // Нулевая сила если есть русские буквы
  }
  
  let strength = 0;
  
  // 1. Длина (макс. 3 балла)
  if (pass.length >= 8) strength += 1;
  if (pass.length >= 12) strength += 1;
  if (pass.length >= 16) strength += 1;

  // 2. Сложность (макс. 5 баллов)
  if (/[A-Z]/.test(pass)) strength += 1;       // +1 заглавные
  if (/[a-z]/.test(pass)) strength += 0.5;     // +0.5 строчные
  if (/\d/.test(pass)) strength += 1;          // +1 цифры
  if (/\d.*\d/.test(pass)) strength += 0.5;    // +0.5 ≥2 цифр
  if (/[!@#$%^&*]/.test(pass)) strength += 1;  // +1 спецсимвол
  if (/[!@#$%^&*].*[!@#$%^&*]/.test(pass)) strength += 1; // +1 ≥2 спецсимволов
 
  // Округляем ВНИЗ до целого числа
  return Math.min(Math.floor(strength), 8);
});

const showPassword = ref(false)

const tooglePassowrdShowing = () => {
  showPassword.value = !showPassword.value
}



const auth = useAuthStore();

const isFocused = ref(false)
const email = ref('');
const emailStatus = ref(0)

const emailReadyToSend = ref(false)
const emailApproved = ref(false)
const emailAnimUpdated = ref(false)

const timerCount = ref(0)
const verifCodeSent = ref(false)
const verifCodeStatus = ref(0)
const canResend = ref(true)
let timerInterval = null

const code = ref('');

const isUNAnimChanged = ref(false)
const username = ref('');
const isTypingUN = ref(false)
let timeoutUN = null

watch(username, async ()=> {
  isTypingUN.value = true
  clearTimeout(timeoutUN)
  timeoutUN = setTimeout(()=> {
    isTypingUN.value = false
    isUNValid.value = validateUN(username.value)
    if(!isUNValid.value && isUNAnimChanged.value) {
      isUNAnimChanged.value = false
      updateCurrentStepFillPercent(-15)
    }
    if (isUNValid.value) {
      validateUNonServer(username.value)
    } else {isUNvalidByServer.value = false}
  }, 1000)
}) 

const validateUNonServer = async(username) => {
  try {
    const {data} = await checkUsername(username)

    if (data?.available) {
      isUNvalidByServer.value = true
      if(!isUNAnimChanged.value) {
        isUNAnimChanged.value = true
        updateCurrentStepFillPercent(15)
      }
    }
    else {
      isUNvalidByServer.value = false
      if(isUNAnimChanged.value) {
        isUNAnimChanged.value = false
        updateCurrentStepFillPercent(-15)
      }
    }
  }
  catch (error) {
    console.error('Ошибка: ', error.response?.data?.message || error.message)
  }
}

const validateUN = (un)=> {
  return /^[a-zA-Z0-9_]{3,20}$/.test(un)
}

let debounceTimer = null

watch(verifCodeStatus, () => {
  const timer = setTimeout(()=> {
    verifCodeStatus.value = 0
    clearTimeout(timer)
  }, 2500)
})

const sendEmailToVerif = () => {
  emailReadyToSend.value = false
  setTimeout(()=> {
    emailApproved.value = true
    sendVerifCode()
  }, 1000)
}

const sendVerifCode = async() => {
  if (!timerInterval) {
    try {
      const {data} = await sendVerificationCodeOnServer(email.value)
      if (data?.status == "success") {
        verifCodeSent.value = true
        startCountdown()
      }
    }
    catch (error) {
      console.error('Ошибка: ', error.response?.data?.message || error.message)
    }
    
  } else notify.warning('Вы уже отправляли код аутентификации!')
  
}

const verifyCode = async () => {
  if(code.value.length == 6) {
    verifCodeStatus.value = -1
    try {
      const {data} = await VerifyCodeOnServer(email.value, code.value)
      if(data?.status == 'success') {
        emailAuthDone.value=true
        verifCodeStatus.value = 1
        updateCurrentStepFillPercent(20)
      }
    }
    catch(error) {
      verifCodeStatus.value = -2
      notify.error(error.response?.data?.message || error.message)
      console.error('Ошибка: ', error.response?.data?.message || error.message)
    }
  } else notify.warning('Неверный синтаксис!')
  
};

const startCountdown = () => {
  if (!timerInterval) {
  timerCount.value = 60

  timerInterval = setInterval(()=> {
    timerCount.value--
    if (timerCount.value <= 0) {
      clearInterval(timerInterval)
      timerInterval = null
      canResend.value=true
      verifCodeSent.value=false
    }
  }, 1000)
  }
} 



const onEmailInput = () => {
  clearTimeout(debounceTimer)
  emailStatus.value=0
  debounceTimer = setTimeout(()=> {
    if (!email.value && emailAnimUpdated.value) {
      updateCurrentStepFillPercent(-20)
      emailAnimUpdated.value = false
    }
    validateEmail()
  }, 1000)
}

const validateEmail = async  () => {
  if (!email.value) {
    return
  }
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if(!emailRegex.test(email.value)) {
    notify.warning('Ошибка синтаксиса...')
    if (emailAnimUpdated.value) {
      updateCurrentStepFillPercent(-20)
      emailAnimUpdated.value = false
    }
    code.value=""
    emailStatus.value = -1
    emailApproved.value=false
    emailReadyToSend.value = false
    return
  }

  try {
    const {data} = await checkEmailOnServer(email.value)
    if (data?.available) {
      emailStatus.value = 1
      if (!emailAnimUpdated.value){
        updateCurrentStepFillPercent(20)
        emailAnimUpdated.value=true
      }
      emailReadyToSend.value = true
    } else if(!data?.available) {
      notify.error('E-mail уже занят!')
      emailStatus.value = -1
      emailApproved.value = false
      emailReadyToSend.value = false
    }
  }
  catch (error) {
    notify.error('Ошибка сервера при проверке E-mail')
    emailStatus.value = -1
    emailApproved.value = false
    emailReadyToSend.value = false
  }
}

const completeRegistration = async () => {
  if (!bottomButtonActive.value) return
  
  if (regStep.value == 1) {
    try {
      const authStore = useAuthStore()

      const data = await authStore.register({
        username: username.value,
        email: email.value,
        password: password.value
      })
      console.log('REG RESULT: ', data)
      if (data?.success) {
        notify.success('Регистрация прошла успешно!')
        closeAuthModal()
      } else {
        notify.error(data.error || "Ошибка регистрации")
      }
    }
    catch(error) {
      notify.error('Неизвестная ошибка при регистрации')
      console.error('Registration error:', error? 'Some error':'Some error')
    }
  } 
};

const closeAuthModal = () => {
  auth.authModalOpen = false;
}

onUnmounted(()=> {
  clearTimeout(timeoutPassword)
  clearTimeout(timeoutPasswordR)
})

</script>
<style scoped>


.auth-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0,0,0,0.9);
  z-index: 800;
  display: flex;
  align-items: center;
  justify-content: center;
}

.auth-modal__content {
  position: relative;
  background: #17212b;
  border-radius: 8px;
  width: 100%;
  height: 500px;
  max-width: 400px;
  animation: border .3s;
}

.auth-content-1 {
  padding: 15px 0px;
  height: 100%;
}

.auth-modal__steps {
    position: relative;
    width: 600px;
    height: 500px;
    border-radius: 8px;
    background-color: #17212b;
    margin-left: 20px;
    background-color: #0e1621;
    overflow: hidden;
}

.auth-header-1 {
    display: flex;
    align-items: center;
    justify-content: center;
    text-align: center;
    font-size: 24px;
    left: 0;
    right: 0;
    height: 80px;
    background-color: #0e1621;
    font-weight: bold;
    box-shadow: 0px 4px 8px black;
    animation: 1s ease-in-out start-auth-header-anim;
}

@keyframes start-auth-header-anim {
  0% {
    height: 250px;
    font-size: 28px;
  }
  100% {
    height: 80px;
    font-size: 24px;
  }
}



.auth-anim {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    top: 35%;
    background-color: #17212b;
    box-shadow: 0px -1px 14px black;
    animation: 1s ease-in-out auth-anim-move;
}

.auth-steps-left {
    display: flex;
    width: 100%;
    align-items: center;
    justify-content: center;
    height: 35%;
}

.auth-step-anim {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #17212b;
    border-radius: 50%;
    margin-left: 30px;
    pointer-events: none;
    user-select: none;
    box-shadow: 2px 2px 8px black;
}

.auth-step-anim::after {
    position: absolute;
    content: "";
    width: 35px;
    height: 10px;
    right: -35px;
    background-color: #17212b;
    box-shadow: 2px 3px 2px rgba(0, 0, 0, 0.61);
}

.auth-step-anim:last-child::after {
    content: none;
}





.input-container {
  display: flex;
  align-items: center;
  border-radius: 8px;
  background-color: #0e1621;
  padding: 5px;
  box-shadow: 2px 2px 4px black;
  margin-right: 15px;
  margin-left: 15px;
  transition: all .3s ease-in-out;
  border: inset 2px #ffffff00;
}

.input-container:focus-within {
  margin-left: 0;
  margin-right: 0;
  border-radius: 0px;
}

.input-container:focus-within > .input-icon {
  opacity: .9;
}

.input-container input {
  flex: 1;
  background-color: transparent;
  color: rgba(255, 255, 255, 0.805);
  border: none;
  margin-left: 10px;
}

.input-container input:-webkit-autofill {
  -webkit-box-shadow: 0 0 0 1000px transparent inset !important; 
  background-color: transparent !important;
  transition: all 5000s ease-in-out 0s !important;
}

.input-container:not(:first-child) {
  margin-top: 20px;
}

.input-container input:focus {
  outline: none;
}

.input-icon {
  opacity: .5;
  width: 32px;
  height: 32px;
  transition: all .3s ease-in-out;
}

.input-container>input:focus ~ .input-icon {
  opacity: .9;
}



.email-sent-btn {
  background-color: #ffffff00;
  transition: all .3s ease-in-out;
  border-radius: 8px;
  cursor: pointer;
  padding: 3px;
}



.email-sent-btn:hover {
  background-color: #ffffff1e;
}

.email-sent-btn svg {
  width: 25px;
  height: 25px;
  opacity: .5;
  scale: .9 1.2;
  rotate: -90deg;
}

.email-menu {
  position: relative;
  display: flex;
  align-items: center;
  height: 100%;
}

.email-status-icon {
  transform: translate(0, 2.5px);
  opacity: 1;
  transition: all .2s ease-in-out;
  margin-right: 5px;
}
.email-verif-code {
  position: absolute;
  display: flex;
  align-items: center;
  height: 46px;
  left: 15px;
  padding: 5px;
  pointer-events: none;
  opacity: 0;
  top: 95px;
  border-radius: 8px;
  background-color: #0e1621;
  box-shadow: 2px 3px 4px black;
  transition: all .4s ease-in-out;
}

.email-verif-code-active {
  opacity: 1;
  top: 150px;
  pointer-events: all;
}

.verify-input {
  background-color: transparent;
  color: rgba(255, 255, 255, 0.805);
  border: none;
  width: 100px;
  margin-left: 15px;
}

.verif-input-icon {
  transition: all .3s ease-in-out;
  margin-right: 5px;
}

.verif-input-icon > svg {
  width: 24px;
  height: 24px;
  cursor: pointer;
  transform: translate(0, 2px);
  opacity: .8;
  transition: all .3s ease-in-out;
}

.verif-input-icon > span {
  opacity: .5;
  pointer-events: none;
  user-select: none;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
}

.verif-input-icon:hover > svg {
  opacity: 1;
}

.send-code-container {
  display: flex;
  align-items: center;
}

.email-verif-code > .input-icon {
  transform: translate(7px, 0);
}

.verify-input:focus {
  outline: none;
}

@keyframes auth-anim-move {
    from {
        top: 0;
    }
    to {
        top: 35%;
    }
}

.email-auth-end {
  pointer-events: none;
  user-select: none;
}

.email-auth-end div:not(.email-status-icon) {
  opacity: .5;
}

.verify-loading {
  transition: all .2s ease-in-out;
  transform: translate(0, 2px);
}

.password {
  position: relative;
}



.password:focus-within .password-strength-container {
  opacity: 1;
  transform: translate(-50%, 100%);
}

.password-strength-container {
  position: absolute;
  left: 50%;
  top: 50%;
  pointer-events: none;
  transform: translate(-50%, -50%);
  border-radius: 8px;
  width: 200px;
  opacity: 0;
  height: 30px;
  background-color: #0a1018;
  transition: all .3s ease;
  box-shadow: 2px 2px 4px black;
  display: flex;
  justify-content: center;
  align-items: center;
}

.password-strength-node {
  background-color: rgb(24, 48, 70);
  border-radius: 50%;
  width: 13px;
  height: 13px;
  box-shadow: 1px 1px 3px black;
  animation: bounce 1s ease-in-out 1s infinite;
  transition: all .3s ease-in-out;
}

.password-strength-node:not(:first-child) {
  margin-left: 9px;
}

.bad-1 {
  animation: bounce 1s ease-in-out 0s infinite;
}
.bad-2 {
  animation: bounce 1s ease-in-out 0.1s infinite;
}
.warning-1 {
  animation: bounce 1s ease-in-out 0.2s infinite;
}

.warning-2 {
  animation: bounce 1s ease-in-out 0.3s infinite;
}

.good-1 {
  animation: bounce 1s ease-in-out 0.4s infinite;
}
.good-2 {
  animation: bounce 1s ease-in-out 0.5s infinite;
}
.awesome-1 {
  animation: bounce 1s ease-in-out 0.6s infinite;
}
.awesome-2 {
  animation: bounce 1s ease-in-out 0.7s infinite;
}

.bad-11 {
  width: 15px;
  height: 15px;
  background-color: #ff1f1f;
  animation: none;
}
.bad-22 {
  width: 15px;
  height: 15px;
  background-color: #ff4800;
  animation: none;
}
.warning-11 {
  animation: none;
  width: 15px;
  height: 15px;
  background-color: #f9fd04;
}
.warning-22 {
  animation: none;
  width: 15px;
  height: 15px;
  background-color: #b3ff00;
}
.good-11 {
  animation: none;
  width: 15px;
  height: 15px;
  background-color: #51ff00;
}
.good-22 {
  animation: none;
  width: 15px;
  height: 15px;
  background-color: #05fc80;
}
.awesome-11 {
  animation: none;
  width: 15px;
  height: 15px;
  background-color: #0099ff;
}
.awesome-22 {
  animation: none;
  width: 15px;
  height: 15px;
  background-color: #0066ff;
}

.input-status {
  position: relative;
  margin-top: 2px;
}

.input-status > svg {
  position: absolute;
  transform: translate(-120%, -50%);
}



.password-show {
  position: relative;
}

.password-show svg {
  position: absolute;
  transform: translate(-260%, -45%);
  opacity: .8;
  cursor: pointer;
  user-select: none;
}

.password-show-icon {
  transition: all .3s ease-in-out;
}

.password-hide-icon {
  transition: all .3s ease-in-out;
}

.bottom-side {
  position: absolute;
  bottom: 0;
  left: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  right: 0;
  height: 80px;
  background-color: #0e1621;
  box-shadow: 0px -4px 8px black;
  animation: 1s ease-in-out start-verif-anim;
  transition: all .3s ease-in-out;
}

.skip-steps-button {
  display: flex;
  text-align: center;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  height: 50px;
  width: 140px;
  margin-left: 30px;
  background-color: #0a1018;
  box-shadow: 2px 2px 6px black;
  cursor: pointer;
  transition: all .3s ease-in-out;
  user-select: none;
  color: #ffffffce;
}

.next-step-button {
  display: flex;
  text-align: center;
  justify-content: center;
  align-items: center;
  border-radius: 8px;
  margin-right: 30px;
  height: 50px;
  width: 180px;
  font-weight: bold;
  background-color: #0a1018;
  box-shadow: 2px 2px 6px rgb(0, 0, 0);
  cursor: pointer;
  transition: all .3s ease-in-out;
  user-select: none;
  color: #ffffffce;
}

.next-step-button:hover {
  background-color: #00ff15;
  color: black;
}

.skip-steps-button:hover {
  background-color: rgb(0, 153, 255);
  color: black;
}

@keyframes start-verif-anim {
  0% {
    height: 50%;
  }
  100% {
    height: 80px;
  }
}



@keyframes bounce {
  0%, 100% {
    transform: translateY(0);
  }
  25% {
    transform: translateY(-2px);
  }
  50% {
    transform: translateY(1px);
  }
  
}

.disabled-btm-btn {
  pointer-events: none;
  color: #ffffff46;
}

.username > svg {
  width: 32px;
  height: 32px;
}

</style>