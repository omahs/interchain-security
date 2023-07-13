"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[815],{3905:(e,t,n)=>{n.d(t,{Zo:()=>p,kt:()=>m});var a=n(7294);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function o(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},i=Object.keys(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var c=a.createContext({}),l=function(e){var t=a.useContext(c),n=t;return e&&(n="function"==typeof e?e(t):o(o({},t),e)),n},p=function(e){var t=l(e.components);return a.createElement(c.Provider,{value:t},e.children)},u="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},h=a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,i=e.originalType,c=e.parentName,p=s(e,["components","mdxType","originalType","parentName"]),u=l(n),h=r,m=u["".concat(c,".").concat(h)]||u[h]||d[h]||i;return n?a.createElement(m,o(o({ref:t},p),{},{components:n})):a.createElement(m,o({ref:t},p))}));function m(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var i=n.length,o=new Array(i);o[0]=h;var s={};for(var c in t)hasOwnProperty.call(t,c)&&(s[c]=t[c]);s.originalType=e,s[u]="string"==typeof e?e:r,o[1]=s;for(var l=2;l<i;l++)o[l]=n[l];return a.createElement.apply(null,o)}return a.createElement.apply(null,n)}h.displayName="MDXCreateElement"},177:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>o,default:()=>d,frontMatter:()=>i,metadata:()=>s,toc:()=>l});var a=n(7462),r=(n(7294),n(3905));const i={sidebar_position:2,title:"Joining Replicated Security testnet"},o=void 0,s={unversionedId:"validators/joining-testnet",id:"validators/joining-testnet",title:"Joining Replicated Security testnet",description:"Introduction",source:"@site/docs/validators/joining-testnet.md",sourceDirName:"validators",slug:"/validators/joining-testnet",permalink:"/interchain-security/validators/joining-testnet",draft:!1,tags:[],version:"current",sidebarPosition:2,frontMatter:{sidebar_position:2,title:"Joining Replicated Security testnet"},sidebar:"tutorialSidebar",previous:{title:"Overview",permalink:"/interchain-security/validators/overview"},next:{title:"Withdrawing consumer chain validator rewards",permalink:"/interchain-security/validators/withdraw_rewards"}},c={},l=[{value:"Introduction",id:"introduction",level:2},{value:"Joining the provider chain",id:"joining-the-provider-chain",level:2},{value:"Initialization",id:"initialization",level:2},{value:"Joining consumer chains",id:"joining-consumer-chains",level:2},{value:"Re-using consensus key",id:"re-using-consensus-key",level:2},{value:"Assigning consensus keys",id:"assigning-consensus-keys",level:2},{value:"Baryon",id:"baryon",level:2},{value:"Noble",id:"noble",level:2}],p={toc:l},u="wrapper";function d(e){let{components:t,...n}=e;return(0,r.kt)(u,(0,a.Z)({},p,n,{components:t,mdxType:"MDXLayout"}),(0,r.kt)("h2",{id:"introduction"},"Introduction"),(0,r.kt)("p",null,"This short guide will teach you how to join the ",(0,r.kt)("a",{parentName:"p",href:"https://github.com/cosmos/testnets/tree/master/replicated-security"},"Replicated Security testnet"),"."),(0,r.kt)("p",null,"The experience gained in the testnet will prepare you for validating interchain secured chains."),(0,r.kt)("admonition",{type:"tip"},(0,r.kt)("p",{parentName:"admonition"},"Provider and consumer chain represent distinct networks and infrastructures operated by the same validator set."),(0,r.kt)("p",{parentName:"admonition"},"For general information about running cosmos-sdk based chains check out the ",(0,r.kt)("a",{parentName:"p",href:"https://hub.cosmos.network/main/validators/validator-setup.html"},"validator basics")," and ",(0,r.kt)("a",{parentName:"p",href:"https://docs.cosmos.network/main/run-node/run-node"},"Running a Node section")," of Cosmos SDK docs")),(0,r.kt)("h2",{id:"joining-the-provider-chain"},"Joining the provider chain"),(0,r.kt)("admonition",{type:"info"},(0,r.kt)("p",{parentName:"admonition"},"At present, all validators of the provider chain must also validate all governance approved consumer chains. The consumer chains cannot have a validator set different than the provider, which means they cannot introduce validators that are not also validating the provider chain.")),(0,r.kt)("p",null,"A comprehensive guide is available ",(0,r.kt)("a",{parentName:"p",href:"https://github.com/cosmos/testnets/tree/master/replicated-security/provider"},"here"),"."),(0,r.kt)("h2",{id:"initialization"},"Initialization"),(0,r.kt)("p",null,"First, initialize your ",(0,r.kt)("inlineCode",{parentName:"p"},"$NODE_HOME")," using the ",(0,r.kt)("inlineCode",{parentName:"p"},"provider")," chain binary."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-bash"},"NODE_MONIKER=<your_node>\nCHAIN_ID=provider\nNODE_HOME=<path_to_your_home>\n\ngaiad init $NODE_MONIKER --chain-id $CHAIN_ID --home $NODE_HOME\n")),(0,r.kt)("p",null,"Add your key to the keyring - more details available ",(0,r.kt)("a",{parentName:"p",href:"https://docs.cosmos.network/main/run-node/keyring"},"here"),"."),(0,r.kt)("p",null,"In this example we will use the ",(0,r.kt)("inlineCode",{parentName:"p"},"test")," keyring-backend. This option is not safe to use in production."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-bash"},"gaiad keys add <key_moniker> --keyring-backend test\n\n# save the address as variable for later use\nMY_VALIDATOR_ADDRESS=$(gaiad keys show my_validator -a --keyring-backend test)\n")),(0,r.kt)("p",null,"Before issuing any transactions, use the ",(0,r.kt)("inlineCode",{parentName:"p"},"provider")," testnet faucet to add funds to your address."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-bash"},'curl https://faucet.rs-testnet.polypore.xyz/request?address=$MY_VALIDATOR_ADDRESS&chain=provider\n\n# example output:\n{\n    "address": "cosmos17p3erf5gv2436fd4vyjwmudakts563a497syuz",\n    "amount": "10000000uatom",\n    "chain": "provider",\n    "hash": "10BFEC53C80C9B649B66549FD88A0B6BCF09E8FCE468A73B4C4243422E724985",\n    "status": "success"\n}\n')),(0,r.kt)("p",null,"Then, use the account associated with the keyring to issue a ",(0,r.kt)("inlineCode",{parentName:"p"},"create-validator")," transaction which will register your validator on chain."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-bash"},'gaiad tx staking create-validator \\\n  --amount=1000000uatom \\\n  --pubkey=$(gaiad tendermint show-validator) \\\n  --moniker="choose a moniker" \\\n  --chain-id=$CHAIN_ID" \\\n  --commission-rate="0.10" \\\n  --commission-max-rate="0.20" \\\n  --commission-max-change-rate="0.01" \\\n  --min-self-delegation="1000000" \\\n  --gas="auto" \\\n  --gas-prices="0.0025uatom" \\\n  --from=<key_moniker>\n')),(0,r.kt)("admonition",{type:"tip"},(0,r.kt)("p",{parentName:"admonition"},"Check this ",(0,r.kt)("a",{parentName:"p",href:"https://hub.cosmos.network/main/validators/validator-setup.html#edit-validator-description"},"guide")," to edit your validator.")),(0,r.kt)("p",null,"After this step, your validator is created and you can start your node and catch up to the rest of the network. It is recommended that you use ",(0,r.kt)("inlineCode",{parentName:"p"},"statesync")," to catch up to the rest of the network."),(0,r.kt)("p",null,"You can use this script to modify your ",(0,r.kt)("inlineCode",{parentName:"p"},"config.toml")," with the required statesync parameters."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-bash"},"# create the statesync script\n$: cd $NODE_HOME\n$: touch statesync.sh\n$ chmod 700 statesync.sh # make executable\n")),(0,r.kt)("p",null,"Paste the following instructions into the ",(0,r.kt)("inlineCode",{parentName:"p"},"statesync.sh"),":"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-bash"},'#!/bin/bash\n\nSNAP_RPC="https://rpc.provider-state-sync-01.rs-testnet.polypore.xyz:443"\n\nLATEST_HEIGHT=$(curl -s $SNAP_RPC/block | jq -r .result.block.header.height); \\\nBLOCK_HEIGHT=$((LATEST_HEIGHT - 2000)); \\\nTRUST_HASH=$(curl -s "$SNAP_RPC/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)\n\nsed -i.bak -E "s|^(enable[[:space:]]+=[[:space:]]+).*$|\\1true| ; \\\ns|^(rpc_servers[[:space:]]+=[[:space:]]+).*$|\\1\\"$SNAP_RPC,$SNAP_RPC\\"| ; \\\ns|^(trust_height[[:space:]]+=[[:space:]]+).*$|\\1$BLOCK_HEIGHT| ; \\\ns|^(trust_hash[[:space:]]+=[[:space:]]+).*$|\\1\\"$TRUST_HASH\\"|" $NODE_HOME/config/config.toml\n')),(0,r.kt)("p",null,"Then, you can execute the script:"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-bash"},"$: ./statesync.sh # setup config.toml for statesync\n")),(0,r.kt)("p",null,"Finally, copy the provider genesis and start your node:"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-bash"},'$: GENESIS_URL=https://github.com/cosmos/testnets/raw/master/replicated-security/provider/provider-genesis.json\n$: wget $GENESIS_URL -O genesis.json\n$: genesis.json $NODE_HOME/config/genesis.json\n# start the service\n$: gaiad start --x-crisis-skip-assert-invariants --home $NODE_HOME --p2p.seeds="08ec17e86dac67b9da70deb20177655495a55407@provider-seed-01.rs-testnet.polypore.xyz:26656,4ea6e56300a2f37b90e58de5ee27d1c9065cf871@provider-seed-02.rs-testnet.polypore.xyz:26656"\n')),(0,r.kt)("p",null,"Additional scripts to setup your nodes are available ",(0,r.kt)("a",{parentName:"p",href:"https://github.com/cosmos/testnets/blob/master/replicated-security/provider/join-rs-provider.sh"},"here")," and ",(0,r.kt)("a",{parentName:"p",href:"https://github.com/cosmos/testnets/blob/master/replicated-security/provider/join-rs-provider-cv.sh"},"here"),". The scripts will configure your node and create the required services - the scripts only work in linux environments."),(0,r.kt)("h2",{id:"joining-consumer-chains"},"Joining consumer chains"),(0,r.kt)("admonition",{type:"tip"},(0,r.kt)("p",{parentName:"admonition"},"Once you reach the active set on the provider chain, you will be required to validate all available consumer chains."),(0,r.kt)("p",{parentName:"admonition"},"You can use the same consensus key on all consumer chains, or opt to use a different key on each consumer chain.\nCheck out this ",(0,r.kt)("a",{parentName:"p",href:"/interchain-security/features/key-assignment"},"guide")," to learn more about key assignment in replicated security.")),(0,r.kt)("p",null,"To join consumer chains, simply replicate the steps above for each consumer using the correct consumer chain binaries."),(0,r.kt)("admonition",{type:"info"},(0,r.kt)("p",{parentName:"admonition"},"When running the provider chain and consumers on the same machine please update the ",(0,r.kt)("inlineCode",{parentName:"p"},"PORT")," numbers for each of them and make sure they do not overlap (otherwise the binaries will not start)."),(0,r.kt)("p",{parentName:"admonition"},"Important ports to re-configure:"),(0,r.kt)("ul",{parentName:"admonition"},(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"--rpc.laddr")),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"--p2p.laddr")),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"--api.address")),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"--grpc.address")),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"--grpc-web.address")))),(0,r.kt)("h2",{id:"re-using-consensus-key"},"Re-using consensus key"),(0,r.kt)("p",null,"To reuse the key on the provider and consumer chains, simply initialize your consumer chain and place the ",(0,r.kt)("inlineCode",{parentName:"p"},"priv_validator_key.json")," into the home directory of your consumer chain (",(0,r.kt)("inlineCode",{parentName:"p"},"<consumer_home>/config/priv_validator_key.json"),")."),(0,r.kt)("p",null,"When you start the chain, the consensus key will be the same on the provider and the consumer chain."),(0,r.kt)("h2",{id:"assigning-consensus-keys"},"Assigning consensus keys"),(0,r.kt)("p",null,"Whenever you initialize a new node, it will be configured with a consensus key you can use."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-bash"},'# machine running consumer chain\nconsumerd init <node_moniker> --home <home_path> --chain-id consumer-1\n\n# use the output of this command to get the consumer chain consensus key\nconsumerd tendermint show-validator\n# output: {"@type":"/cosmos.crypto.ed25519.PubKey","key":"<key>"}\n')),(0,r.kt)("p",null,"Then, let the provider know which key you will be using for the consumer chain:"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-bash"},"# machine running the provider chain\ngaiad tx provider assign-consensus-key consumer-1 '<consumer_pubkey>' --from <key_moniker> --home $NODE_HOME --gas 900000 -b sync -y -o json\n")),(0,r.kt)("p",null,"After this step, you are ready to copy the consumer genesis into your nodes's ",(0,r.kt)("inlineCode",{parentName:"p"},"/config")," folder, start your consumer chain node and catch up to the network."),(0,r.kt)("h2",{id:"baryon"},"Baryon"),(0,r.kt)("p",null,"You can find the onboarding repo instructions for the Baryon chain ",(0,r.kt)("a",{parentName:"p",href:"https://github.com/cosmos/testnets/blob/master/replicated-security/stopped/baryon-1/README.md"},"here")),(0,r.kt)("h2",{id:"noble"},"Noble"),(0,r.kt)("p",null,"You can find the onboarding repo instructions for the Noble chain ",(0,r.kt)("a",{parentName:"p",href:"https://github.com/cosmos/testnets/blob/master/replicated-security/stopped/noble-1/README.md"},"here")))}d.isMDXComponent=!0}}]);