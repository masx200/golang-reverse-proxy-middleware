if(!self.define){let e,n={};const o=(o,s)=>(o=new URL(o+".js",s).href,n[o]||new Promise((n=>{if("document"in self){const e=document.createElement("script");e.src=o,e.onload=n,document.head.appendChild(e)}else e=o,importScripts(o),n()})).then((()=>{let e=n[o];if(!e)throw new Error(`Module ${o} didn’t register its module`);return e})));self.define=(s,t)=>{const i=e||("document"in self?document.currentScript.src:"")||location.href;if(n[i])return;let r={};const l=e=>o(e,i),u={module:{uri:i},exports:r,require:l};n[i]=Promise.all(s.map((e=>u[e]||l(e)))).then((e=>(t(...e),r)))}}define(["./workbox-f3e6b16a"],(function(e){"use strict";self.skipWaiting(),e.clientsClaim(),e.precacheAndRoute([{url:"assets/background_结果-C6o_kEWo.webp",revision:null},{url:"assets/favicon_结果-HVYl5oM6.webp",revision:null},{url:"manifest.webmanifest",revision:"5c55f282eb99ac7d89d3652973a125f8"}],{}),e.cleanupOutdatedCaches(),e.registerRoute(new e.NavigationRoute(e.createHandlerBoundToURL("index.html")))}));