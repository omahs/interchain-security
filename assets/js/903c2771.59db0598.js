"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[54],{3905:(e,t,o)=>{o.d(t,{Zo:()=>c,kt:()=>m});var n=o(7294);function a(e,t,o){return t in e?Object.defineProperty(e,t,{value:o,enumerable:!0,configurable:!0,writable:!0}):e[t]=o,e}function r(e,t){var o=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),o.push.apply(o,n)}return o}function i(e){for(var t=1;t<arguments.length;t++){var o=null!=arguments[t]?arguments[t]:{};t%2?r(Object(o),!0).forEach((function(t){a(e,t,o[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(o)):r(Object(o)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(o,t))}))}return e}function l(e,t){if(null==e)return{};var o,n,a=function(e,t){if(null==e)return{};var o,n,a={},r=Object.keys(e);for(n=0;n<r.length;n++)o=r[n],t.indexOf(o)>=0||(a[o]=e[o]);return a}(e,t);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);for(n=0;n<r.length;n++)o=r[n],t.indexOf(o)>=0||Object.prototype.propertyIsEnumerable.call(e,o)&&(a[o]=e[o])}return a}var s=n.createContext({}),u=function(e){var t=n.useContext(s),o=t;return e&&(o="function"==typeof e?e(t):i(i({},t),e)),o},c=function(e){var t=u(e.components);return n.createElement(s.Provider,{value:t},e.children)},p="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},h=n.forwardRef((function(e,t){var o=e.components,a=e.mdxType,r=e.originalType,s=e.parentName,c=l(e,["components","mdxType","originalType","parentName"]),p=u(o),h=a,m=p["".concat(s,".").concat(h)]||p[h]||d[h]||r;return o?n.createElement(m,i(i({ref:t},c),{},{components:o})):n.createElement(m,i({ref:t},c))}));function m(e,t){var o=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var r=o.length,i=new Array(r);i[0]=h;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l[p]="string"==typeof e?e:a,i[1]=l;for(var u=2;u<r;u++)i[u]=o[u];return n.createElement.apply(null,i)}return n.createElement.apply(null,o)}h.displayName="MDXCreateElement"},8896:(e,t,o)=>{o.r(t),o.d(t,{assets:()=>s,contentTitle:()=>i,default:()=>d,frontMatter:()=>r,metadata:()=>l,toc:()=>u});var n=o(7462),a=(o(7294),o(3905));const r={sidebar_position:10,title:"Soft Opt-Out"},i=void 0,l={unversionedId:"adrs/adr-009-soft-opt-out",id:"adrs/adr-009-soft-opt-out",title:"Soft Opt-Out",description:"ADR 009: Soft Opt-Out",source:"@site/docs/adrs/adr-009-soft-opt-out.md",sourceDirName:"adrs",slug:"/adrs/adr-009-soft-opt-out",permalink:"/interchain-security/adrs/adr-009-soft-opt-out",draft:!1,tags:[],version:"current",sidebarPosition:10,frontMatter:{sidebar_position:10,title:"Soft Opt-Out"},sidebar:"tutorialSidebar",previous:{title:"Throttle with retries",permalink:"/interchain-security/adrs/adr-008-throttle-retries"}},s={},u=[{value:"ADR 009: Soft Opt-Out",id:"adr-009-soft-opt-out",level:2},{value:"Changelog",id:"changelog",level:2},{value:"Status",id:"status",level:2},{value:"Context",id:"context",level:2},{value:"Decision",id:"decision",level:2},{value:"Consequences",id:"consequences",level:2},{value:"Positive",id:"positive",level:3},{value:"Negative",id:"negative",level:3},{value:"Neutral",id:"neutral",level:3},{value:"References",id:"references",level:2}],c={toc:u},p="wrapper";function d(e){let{components:t,...o}=e;return(0,a.kt)(p,(0,n.Z)({},c,o,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h2",{id:"adr-009-soft-opt-out"},"ADR 009: Soft Opt-Out"),(0,a.kt)("h2",{id:"changelog"},"Changelog"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"6/13/23: Initial draft of ADR. Feature already implemented and in production.")),(0,a.kt)("h2",{id:"status"},"Status"),(0,a.kt)("p",null,"Accepted"),(0,a.kt)("h2",{id:"context"},"Context"),(0,a.kt)("p",null,"Some small validators may not have the resources needed to validate all consumer chains. Therefore a need exists to allow the bottom ",(0,a.kt)("inlineCode",{parentName:"p"},"x%")," of validators to opt-out of validating a consumer chain. Meaning downtime infractions for these validators are dropped without ever reaching the provider."),(0,a.kt)("p",null,"This document specifies a modification to the ccv protocol which allows the bottom x% of the validator set by power to opt out of validating consumer chains without being jailed or otherwise punished for it. The feature is implemented with entirely consumer-side code."),(0,a.kt)("h2",{id:"decision"},"Decision"),(0,a.kt)("p",null,"A consumer param exists, known as ",(0,a.kt)("inlineCode",{parentName:"p"},"SoftOptOutThreshold"),", which is a string decimal in the range of ","[0, 0.2]",", that determines the portion of validators which are allowed to opt out of validating that specific consumer."),(0,a.kt)("p",null,"In every consumer beginblocker, a function is ran which determines the so called  ",(0,a.kt)("em",{parentName:"p"},"smallest non opt-out voting power"),". Validators with voting power greater than or equal to this value must validate the consumer chain, while validators below this value may opt out of validating the consumer chain."),(0,a.kt)("p",null,"The smallest non opt-out voting power is recomputed every beginblocker in ",(0,a.kt)("inlineCode",{parentName:"p"},"UpdateSmallestNonOptOutPower()"),". In a nutshell, the method obtains the total voting power of the consumer, iterates through the full valset (ordered power ascending) keeping track of a power sum, and when ",(0,a.kt)("inlineCode",{parentName:"p"},"powerSum / totalPower > SoftOptOutThreshold"),", the ",(0,a.kt)("inlineCode",{parentName:"p"},"SmallestNonOptOutPower")," is found and persisted."),(0,a.kt)("p",null,"Then, whenever the ",(0,a.kt)("inlineCode",{parentName:"p"},"Slash()")," interface is executed on the consumer, if the voting power of the relevant validator being slashed is less than ",(0,a.kt)("inlineCode",{parentName:"p"},"SmallestNonOptOutPower")," for that block, the slash request is dropped and never sent to the provider."),(0,a.kt)("h2",{id:"consequences"},"Consequences"),(0,a.kt)("h3",{id:"positive"},"Positive"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"Small validators can opt out of validating specific consumers without being punished for it.")),(0,a.kt)("h3",{id:"negative"},"Negative"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"The bottom ",(0,a.kt)("inlineCode",{parentName:"li"},"x%")," is still part of the total voting power of the consumer chain. This means that if the soft opt-out threshold is set to ",(0,a.kt)("inlineCode",{parentName:"li"},"10%")," for example, and every validator in the bottom ",(0,a.kt)("inlineCode",{parentName:"li"},"10%")," opts out from validating the consumer, then a ",(0,a.kt)("inlineCode",{parentName:"li"},"24%")," downtime of the remaining voting power would halt the chain. This may be especially problematic during consumer upgrades."),(0,a.kt)("li",{parentName:"ul"},"In nominal scenarios, consumers with soft opt out enabled will be constructing slash packets for small vals, which may be dropped. This is wasted computation, but necessary to keep implementation simple. Note that the sdk's ",(0,a.kt)("a",{parentName:"li",href:"https://github.com/cosmos/cosmos-sdk/blob/d3f09c222243bb3da3464969f0366330dcb977a8/x/slashing/keeper/infractions.go#L75"},"full downtime logic")," is always executed on the consumer, which can be computationally expensive and slow down certain blocks.")),(0,a.kt)("h3",{id:"neutral"},"Neutral"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"Validators in the bottom of the valset who don't have to validate, may receive large delegation(s) which suddenly boost the validator to the subset that has to validate. This may catch the validator off guard.")),(0,a.kt)("h2",{id:"references"},"References"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"Original issue with some napkin math ",(0,a.kt)("a",{parentName:"li",href:"https://github.com/cosmos/interchain-security/issues/784"},"#784"))))}d.isMDXComponent=!0}}]);