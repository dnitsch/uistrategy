package uistrategy_test

var testHtml_style = []byte(`<html lang="en"><head><style>.ͼ1.cm-editor.cm-focused {outline: 1px dotted #212121;}
.ͼ1.cm-editor {position: relative !important; box-sizing: border-box; display: flex !important; flex-direction: column;}
.ͼ1 .cm-scroller {display: flex !important; align-items: flex-start !important; font-family: monospace; line-height: 1.4; height: 100%; overflow-x: auto; position: relative; z-index: 0;}
.ͼ1 .cm-content[contenteditable=true] {-webkit-user-modify: read-write-plaintext-only;}
.ͼ1 .cm-content {margin: 0; flex-grow: 2; flex-shrink: 0; min-height: 100%; display: block; white-space: pre; word-wrap: normal; box-sizing: border-box; padding: 4px 0; outline: none;}
.ͼ1 .cm-lineWrapping {white-space: pre-wrap; white-space: break-spaces; word-break: break-word; overflow-wrap: anywhere; flex-shrink: 1;}
.ͼ2 .cm-content {caret-color: black;}
.ͼ3 .cm-content {caret-color: white;}
.ͼ1 .cm-line {display: block; padding: 0 2px 0 4px;}
.ͼ1 .cm-selectionLayer {z-index: -1; contain: size style;}
.ͼ1 .cm-selectionBackground {position: absolute;}
.ͼ2 .cm-selectionBackground {background: #d9d9d9;}
.ͼ3 .cm-selectionBackground {background: #222;}
.ͼ2.cm-focused .cm-selectionBackground {background: #d7d4f0;}
.ͼ3.cm-focused .cm-selectionBackground {background: #233;}
.ͼ1 .cm-cursorLayer {z-index: 100; contain: size style; pointer-events: none;}
.ͼ1.cm-focused .cm-cursorLayer {animation: steps(1) cm-blink 1.2s infinite;}
@keyframes cm-blink {50% {opacity: 0;}}
@keyframes cm-blink2 {50% {opacity: 0;}}
.ͼ1 .cm-cursor, .ͼ1 .cm-dropCursor {position: absolute; border-left: 1.2px solid black; margin-left: -0.6px; pointer-events: none;}
.ͼ1 .cm-cursor {display: none;}
.ͼ3 .cm-cursor {border-left-color: #444;}
.ͼ1.cm-focused .cm-cursor {display: block;}
.ͼ2 .cm-activeLine {background-color: #cceeff44;}
.ͼ3 .cm-activeLine {background-color: #99eeff33;}
.ͼ2 .cm-specialChar {color: red;}
.ͼ3 .cm-specialChar {color: #f78;}
.ͼ1 .cm-gutters {flex-shrink: 0; display: flex; height: 100%; box-sizing: border-box; left: 0; z-index: 200;}
.ͼ2 .cm-gutters {background-color: #f5f5f5; color: #6c6c6c; border-right: 1px solid #ddd;}
.ͼ3 .cm-gutters {background-color: #333338; color: #ccc;}
.ͼ1 .cm-gutter {display: flex !important; flex-direction: column; flex-shrink: 0; box-sizing: border-box; min-height: 100%; overflow: hidden;}
.ͼ1 .cm-gutterElement {box-sizing: border-box;}
.ͼ1 .cm-lineNumbers .cm-gutterElement {padding: 0 3px 0 5px; min-width: 20px; text-align: right; white-space: nowrap;}
.ͼ2 .cm-activeLineGutter {background-color: #e2f2ff;}
.ͼ3 .cm-activeLineGutter {background-color: #222227;}
.ͼ1 .cm-panels {box-sizing: border-box; position: sticky; left: 0; right: 0;}
.ͼ2 .cm-panels {background-color: #f5f5f5; color: black;}
.ͼ2 .cm-panels-top {border-bottom: 1px solid #ddd;}
.ͼ2 .cm-panels-bottom {border-top: 1px solid #ddd;}
.ͼ3 .cm-panels {background-color: #333338; color: white;}
.ͼ1 .cm-tab {display: inline-block; overflow: hidden; vertical-align: bottom;}
.ͼ1 .cm-widgetBuffer {vertical-align: text-top; height: 1em; width: 0; display: inline;}
.ͼ1 .cm-placeholder {color: #888; display: inline-block; vertical-align: top;}
.ͼ1 .cm-button {vertical-align: middle; color: inherit; font-size: 70%; padding: .2em 1em; border-radius: 1px;}
.ͼ2 .cm-button:active {background-image: linear-gradient(#b4b4b4, #d0d3d6);}
.ͼ2 .cm-button {background-image: linear-gradient(#eff1f5, #d9d9df); border: 1px solid #888;}
.ͼ3 .cm-button:active {background-image: linear-gradient(#111, #333);}
.ͼ3 .cm-button {background-image: linear-gradient(#393939, #111); border: 1px solid #888;}
.ͼ1 .cm-textfield {vertical-align: middle; color: inherit; font-size: 70%; border: 1px solid silver; padding: .2em .5em;}
.ͼ2 .cm-textfield {background-color: white;}
.ͼ3 .cm-textfield {border: 1px solid #555; background-color: inherit;}
.ͼ1 .cm-tooltip.cm-tooltip-autocomplete > ul > li {overflow-x: hidden; text-overflow: ellipsis; cursor: pointer; padding: 1px 3px; line-height: 1.2;}
.ͼ1 .cm-tooltip.cm-tooltip-autocomplete > ul {font-family: monospace; white-space: nowrap; overflow: hidden auto; max-width: 700px; max-width: min(700px, 95vw); min-width: 250px; max-height: 10em; list-style: none; margin: 0; padding: 0;}
.ͼ2 .cm-tooltip-autocomplete ul li[aria-selected] {background: #17c; color: white;}
.ͼ3 .cm-tooltip-autocomplete ul li[aria-selected] {background: #347; color: white;}
.ͼ1 .cm-completionListIncompleteTop:before, .ͼ1 .cm-completionListIncompleteBottom:after {content: "···"; opacity: 0.5; display: block; text-align: center;}
.ͼ1 .cm-tooltip.cm-completionInfo {position: absolute; padding: 3px 9px; width: max-content; max-width: 400px; box-sizing: border-box;}
.ͼ1 .cm-completionInfo.cm-completionInfo-left {right: 100%;}
.ͼ1 .cm-completionInfo.cm-completionInfo-right {left: 100%;}
.ͼ1 .cm-completionInfo.cm-completionInfo-left-narrow {right: 30px;}
.ͼ1 .cm-completionInfo.cm-completionInfo-right-narrow {left: 30px;}
.ͼ2 .cm-snippetField {background-color: #00000022;}
.ͼ3 .cm-snippetField {background-color: #ffffff22;}
.ͼ1 .cm-snippetFieldPosition {vertical-align: text-top; width: 0; height: 1.15em; display: inline-block; margin: 0 -0.7px -.7em; border-left: 1.4px dotted #888;}
.ͼ1 .cm-completionMatchedText {text-decoration: underline;}
.ͼ1 .cm-completionDetail {margin-left: 0.5em; font-style: italic;}
.ͼ1 .cm-completionIcon {font-size: 90%; width: .8em; display: inline-block; text-align: center; padding-right: .6em; opacity: 0.6;}
.ͼ1 .cm-completionIcon-function:after, .ͼ1 .cm-completionIcon-method:after {content: 'ƒ';} 
.ͼ1 .cm-completionIcon-text:after {content: 'abc'; font-size: 50%; vertical-align: middle;}
.ͼ1 .cm-tooltip {z-index: 100;}
.ͼ2 .cm-tooltip {border: 1px solid #bbb; background-color: #f5f5f5;}
.ͼ2 .cm-tooltip-section:not(:first-child) {border-top: 1px solid #bbb;}
.ͼ3 .cm-tooltip {background-color: #333338; color: white;}
.ͼ1 .cm-tooltip-arrow:before, .ͼ1 .cm-tooltip-arrow:after {content: ''; position: absolute; width: 0; height: 0; border-left: 7px solid transparent; border-right: 7px solid transparent;}
.ͼ1 .cm-tooltip-above .cm-tooltip-arrow:before {border-top: 7px solid #bbb;}
.ͼ1 .cm-tooltip-above .cm-tooltip-arrow:after {border-top: 7px solid #f5f5f5; bottom: 1px;}
.ͼ1 .cm-tooltip-above .cm-tooltip-arrow {bottom: -7px;}
.ͼ1 .cm-tooltip-below .cm-tooltip-arrow:before {border-bottom: 7px solid #bbb;}
.ͼ1 .cm-tooltip-below .cm-tooltip-arrow:after {border-bottom: 7px solid #f5f5f5; top: 1px;}
.ͼ1 .cm-tooltip-below .cm-tooltip-arrow {top: -7px;}
.ͼ1 .cm-tooltip-arrow {height: 7px; width: 14px; position: absolute; z-index: -1; overflow: hidden;}
.ͼ3 .cm-tooltip .cm-tooltip-arrow:before {border-top-color: #333338; border-bottom-color: #333338;}
.ͼ3 .cm-tooltip .cm-tooltip-arrow:after {border-top-color: transparent; border-bottom-color: transparent;}
.ͼ1 .cm-selectionMatch {background-color: #99ff7780;}
.ͼ1 .cm-searchMatch .cm-selectionMatch {background-color: transparent;}
.ͼ1.cm-focused .cm-matchingBracket {background-color: #328c8252;}
.ͼ1.cm-focused .cm-nonmatchingBracket {background-color: #bb555544;}
.ͼ5 {color: #7a757a;}
.ͼ6 {text-decoration: underline;}
.ͼ7 {text-decoration: underline; font-weight: bold;}
.ͼ8 {font-style: italic;}
.ͼ9 {font-weight: bold;}
.ͼa {text-decoration: line-through;}
.ͼb {color: #708;}
.ͼc {color: #219;}
.ͼd {color: #164;}
.ͼe {color: #a11;}
.ͼf {color: #e40;}
.ͼg {color: #00f;}
.ͼh {color: #30a;}
.ͼi {color: #085;}
.ͼj {color: #167;}
.ͼk {color: #256;}
.ͼl {color: #00c;}
.ͼm {color: #940;}
.ͼn {color: #f00;}
.ͼ4 .cm-line ::selection {background-color: transparent !important;}
.ͼ4 .cm-line::selection {background-color: transparent !important;}
.ͼ4 .cm-line {caret-color: transparent !important;}
</style>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="Content-Security-Policy" content="default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' http://127.0.0.1:* data:; connect-src 'self' http://127.0.0.1:*; script-src 'self' 'sha256-GRUzBA7PzKYug7pqxv5rJaec5bwDCw1Vo6/IXwvD3Tc='">

    <title>Collections - Acme - PocketBase</title>

    <link rel="apple-touch-icon" sizes="180x180" href="./images/favicon/apple-touch-icon.png">
    <link rel="icon" type="image/png" sizes="32x32" href="./images/favicon/favicon-32x32.png">
    <link rel="icon" type="image/png" sizes="16x16" href="./images/favicon/favicon-16x16.png">
    <link rel="manifest" href="./images/favicon/site.webmanifest">
    <link rel="mask-icon" href="./images/favicon/safari-pinned-tab.svg" color="#000000">
    <link rel="shortcut icon" href="./images/favicon/favicon.ico">
    <meta name="msapplication-TileColor" content="#ffffff">
    <meta name="msapplication-config" content="./images/favicon/browserconfig.xml">
    <meta name="theme-color" content="#ffffff">

    <script>
        window.Prism = window.Prism || {};
        window.Prism.manual = true;
    </script>
  <script type="module" crossorigin="" src="./assets/index.e13041a6.js"></script>
  <link rel="stylesheet" href="./assets/index.26987507.css">
<link rel="modulepreload" as="script" crossorigin="" href="http://localhost:8090/_/assets/FilterAutocompleteInput.37739e76.js"><link rel="modulepreload" as="script" crossorigin="" href="http://localhost:8090/_/assets/index.a9121ab1.js"></head>
<body>
    <div id="app"> <div class="app-layout"><aside class="app-sidebar"><a href="#/" class="logo logo-sm"><img src="./images/logo.svg" alt="PocketBase logo" width="40" height="40"></a> <nav class="main-menu"><a href="#/collections" class="menu-item current-route" aria-label="Collections"><i class="ri-database-2-line"></i></a> <a href="#/users" class="menu-item" aria-label="Users"><i class="ri-group-line"></i></a> <a href="#/logs" class="menu-item" aria-label="Logs"><i class="ri-line-chart-line"></i></a> <a href="#/settings" class="menu-item" aria-label="Settings"><i class="ri-tools-line"></i></a></nav> <figure class="thumb thumb-circle link-hint closable"><img src="./images/avatars/avatar0.svg" alt="Avatar"> <div class="toggler-container"></div></figure></aside> <div class="app-body"><aside class="page-sidebar collection-sidebar"><header class="sidebar-header"><div class="form-field search"><div class="form-field-addon"><button type="button" class="btn btn-xs btn-secondary btn-circle btn-clear hidden"><i class="ri-close-line"></i></button></div> <input type="text" placeholder="Search collections..."></div></header> <hr class="m-t-5 m-b-xs"> <div class="sidebar-content"><div tabindex="0" class="sidebar-list-item active"><i class="ri-folder-open-line"></i> <span class="txt">test</span> </div></div> <footer class="sidebar-footer"><button type="button" class="btn btn-block btn-outline"><i class="ri-add-line"></i> 
                <span class="txt">New collection</span></button></footer></aside>   <div class="page-wrapper"><main class="page-content"><header class="page-header"><nav class="breadcrumbs"><div class="breadcrumb-item">Collections</div> <div class="breadcrumb-item">test</div></nav> <div class="inline-flex gap-5"><button type="button" class="btn btn-secondary btn-circle"><i class="ri-settings-4-line"></i></button> <button type="button" class="btn btn-secondary btn-circle svelte-b7gb6q"><i class="ri-refresh-line svelte-b7gb6q"></i></button></div> <div class="btns-group"><button type="button" class="btn btn-outline"><i class="ri-code-s-slash-line"></i> 
                    <span class="txt">API Preview</span></button> <button type="button" class="btn btn-expanded"><i class="ri-add-line"></i> 
                    <span class="txt">New record</span></button></div></header> <div class="searchbar-wrapper"><form class="searchbar"><label for="search_4Ff0yc2" class="m-l-10 txt-xl"><i class="ri-search-line"></i></label> <div class="code-editor"><div class="cm-editor ͼ1 ͼ2 ͼ4"><div aria-live="polite" style="position: absolute; top: -10000px;"></div><div tabindex="-1" class="cm-scroller"><div spellcheck="false" autocorrect="off" autocapitalize="off" translate="no" contenteditable="true" class="cm-content cm-lineWrapping" style="tab-size: 4;" role="textbox" aria-multiline="true" aria-autocomplete="list"><div class="cm-line"><img class="cm-widgetBuffer" aria-hidden="true"><span class="cm-placeholder" aria-label="placeholder Search filter, ex. created > &quot;2022-01-01&quot;..." contenteditable="false" style="pointer-events: none;">Search filter, ex. created &gt; "2022-01-01"...</span><br></div></div><div class="cm-selectionLayer" aria-hidden="true"></div><div class="cm-cursorLayer" aria-hidden="true" style="animation-duration: 1200ms;"><div class="cm-cursor cm-cursor-primary" style="left: 0px; top: 6.5px; height: 16.5px;"></div></div></div></div></div> </form></div> <div class="table-wrapper"><table class="table"><thead><tr><th class="bulk-select-col min-width"><div class="form-field"><input type="checkbox" id="checkbox_0" disabled=""> <label for="checkbox_0"></label></div></th> <th tabindex="0" class="col-sort col-type-text col-field-id"><div class="col-header-content"><i class="ri-key-line"></i> <span class="txt">id</span></div></th> <th tabindex="0" class="col-sort col-type-text col-field-testField1"><div class="col-header-content"><i class="ri-text"></i> <span class="txt">testField1</span></div></th> <th tabindex="0" class="col-sort col-type-date col-field-created sort-active sort-desc"><div class="col-header-content"><i class="ri-calendar-line"></i> <span class="txt">created</span></div></th> <th tabindex="0" class="col-sort col-type-date col-field-updated"><div class="col-header-content"><i class="ri-calendar-line"></i> <span class="txt">updated</span></div></th> <th class="col-type-action min-width"></th></tr></thead> <tbody><tr><td colspan="99" class="txt-center txt-hint p-xs"><h6>No records found.</h6> </td> </tr></tbody></table></div>   </main> <footer class="page-footer"><a href="https://github.com/pocketbase/pocketbase/releases" class="inline-flex flex-gap-5" target="_blank" rel="noopener" title="Releases"><span class="txt">PocketBase v0.7.9</span></a></footer></div>     <div class="toasts-wrapper"></div></div></div> </div>
    


<div class="overlays"><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div></div><div class="app-tooltip"></div></body></html>`)

var testHtml_noStyle = []byte(`<html lang="en"><head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="Content-Security-Policy" content="default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' http://127.0.0.1:* data:; connect-src 'self' http://127.0.0.1:*; script-src 'self' 'sha256-GRUzBA7PzKYug7pqxv5rJaec5bwDCw1Vo6/IXwvD3Tc='">

    <title>Collections - Acme - PocketBase</title>

    <link rel="apple-touch-icon" sizes="180x180" href="./images/favicon/apple-touch-icon.png">
    <link rel="icon" type="image/png" sizes="32x32" href="./images/favicon/favicon-32x32.png">
    <link rel="icon" type="image/png" sizes="16x16" href="./images/favicon/favicon-16x16.png">
    <link rel="manifest" href="./images/favicon/site.webmanifest">
    <link rel="mask-icon" href="./images/favicon/safari-pinned-tab.svg" color="#000000">
    <link rel="shortcut icon" href="./images/favicon/favicon.ico">
    <meta name="msapplication-TileColor" content="#ffffff">
    <meta name="msapplication-config" content="./images/favicon/browserconfig.xml">
    <meta name="theme-color" content="#ffffff">

    <script>
        window.Prism = window.Prism || {};
        window.Prism.manual = true;
    </script>
  <script type="module" crossorigin="" src="./assets/index.e13041a6.js"></script>
  <link rel="stylesheet" href="./assets/index.26987507.css">
<link rel="modulepreload" as="script" crossorigin="" href="http://localhost:8090/_/assets/FilterAutocompleteInput.37739e76.js"><link rel="modulepreload" as="script" crossorigin="" href="http://localhost:8090/_/assets/index.a9121ab1.js"></head>
<body>
    <div id="app"> <div class="app-layout"><aside class="app-sidebar"><a href="#/" class="logo logo-sm"><img src="./images/logo.svg" alt="PocketBase logo" width="40" height="40"></a> <nav class="main-menu"><a href="#/collections" class="menu-item current-route" aria-label="Collections"><i class="ri-database-2-line"></i></a> <a href="#/users" class="menu-item" aria-label="Users"><i class="ri-group-line"></i></a> <a href="#/logs" class="menu-item" aria-label="Logs"><i class="ri-line-chart-line"></i></a> <a href="#/settings" class="menu-item" aria-label="Settings"><i class="ri-tools-line"></i></a></nav> <figure class="thumb thumb-circle link-hint closable"><img src="./images/avatars/avatar0.svg" alt="Avatar"> <div class="toggler-container"></div></figure></aside> <div class="app-body"><aside class="page-sidebar collection-sidebar"><header class="sidebar-header"><div class="form-field search"><div class="form-field-addon"><button type="button" class="btn btn-xs btn-secondary btn-circle btn-clear hidden"><i class="ri-close-line"></i></button></div> <input type="text" placeholder="Search collections..."></div></header> <hr class="m-t-5 m-b-xs"> <div class="sidebar-content"><div tabindex="0" class="sidebar-list-item active"><i class="ri-folder-open-line"></i> <span class="txt">test</span> </div></div> <footer class="sidebar-footer"><button type="button" class="btn btn-block btn-outline"><i class="ri-add-line"></i> 
                <span class="txt">New collection</span></button></footer></aside>   <div class="page-wrapper"><main class="page-content"><header class="page-header"><nav class="breadcrumbs"><div class="breadcrumb-item">Collections</div> <div class="breadcrumb-item">test</div></nav> <div class="inline-flex gap-5"><button type="button" class="btn btn-secondary btn-circle"><i class="ri-settings-4-line"></i></button> <button type="button" class="btn btn-secondary btn-circle svelte-b7gb6q"><i class="ri-refresh-line svelte-b7gb6q"></i></button></div> <div class="btns-group"><button type="button" class="btn btn-outline"><i class="ri-code-s-slash-line"></i> 
                    <span class="txt">API Preview</span></button> <button type="button" class="btn btn-expanded"><i class="ri-add-line"></i> 
                    <span class="txt">New record</span></button></div></header> <div class="searchbar-wrapper"><form class="searchbar"><label for="search_4Ff0yc2" class="m-l-10 txt-xl"><i class="ri-search-line"></i></label> <div class="code-editor"><div class="cm-editor ͼ1 ͼ2 ͼ4"><div aria-live="polite" style="position: absolute; top: -10000px;"></div><div tabindex="-1" class="cm-scroller"><div spellcheck="false" autocorrect="off" autocapitalize="off" translate="no" contenteditable="true" class="cm-content cm-lineWrapping" style="tab-size: 4;" role="textbox" aria-multiline="true" aria-autocomplete="list"><div class="cm-line"><img class="cm-widgetBuffer" aria-hidden="true"><span class="cm-placeholder" aria-label="placeholder Search filter, ex. created > &quot;2022-01-01&quot;..." contenteditable="false" style="pointer-events: none;">Search filter, ex. created &gt; "2022-01-01"...</span><br></div></div><div class="cm-selectionLayer" aria-hidden="true"></div><div class="cm-cursorLayer" aria-hidden="true" style="animation-duration: 1200ms;"><div class="cm-cursor cm-cursor-primary" style="left: 0px; top: 6.5px; height: 16.5px;"></div></div></div></div></div> </form></div> <div class="table-wrapper"><table class="table"><thead><tr><th class="bulk-select-col min-width"><div class="form-field"><input type="checkbox" id="checkbox_0" disabled=""> <label for="checkbox_0"></label></div></th> <th tabindex="0" class="col-sort col-type-text col-field-id"><div class="col-header-content"><i class="ri-key-line"></i> <span class="txt">id</span></div></th> <th tabindex="0" class="col-sort col-type-text col-field-testField1"><div class="col-header-content"><i class="ri-text"></i> <span class="txt">testField1</span></div></th> <th tabindex="0" class="col-sort col-type-date col-field-created sort-active sort-desc"><div class="col-header-content"><i class="ri-calendar-line"></i> <span class="txt">created</span></div></th> <th tabindex="0" class="col-sort col-type-date col-field-updated"><div class="col-header-content"><i class="ri-calendar-line"></i> <span class="txt">updated</span></div></th> <th class="col-type-action min-width"></th></tr></thead> <tbody><tr><td colspan="99" class="txt-center txt-hint p-xs"><h6>No records found.</h6> </td> </tr></tbody></table></div>   </main> <footer class="page-footer"><a href="https://github.com/pocketbase/pocketbase/releases" class="inline-flex flex-gap-5" target="_blank" rel="noopener" title="Releases"><span class="txt">PocketBase v0.7.9</span></a></footer></div>     <div class="toasts-wrapper"></div></div></div> </div>
    


<div class="overlays"><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div></div><div class="app-tooltip"></div></body></html>`)
