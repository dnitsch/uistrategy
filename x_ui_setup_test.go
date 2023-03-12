package uistrategy_test

var pocketBaseStyle = []byte(`<html lang="en"><body>
    <div id="app"> <div class="app-layout"><aside class="app-sidebar"><a href="#/" class="logo logo-sm"><img src="./images/logo.svg" alt="PocketBase logo" width="40" height="40"></a> <nav class="main-menu"><a href="#/collections" class="menu-item current-route" aria-label="Collections"><i class="ri-database-2-line"></i></a> <a href="#/users" class="menu-item" aria-label="Users"><i class="ri-group-line"></i></a> <a href="#/logs" class="menu-item" aria-label="Logs"><i class="ri-line-chart-line"></i></a> <a href="#/settings" class="menu-item" aria-label="Settings"><i class="ri-tools-line"></i></a></nav> <figure class="thumb thumb-circle link-hint closable"><img src="./images/avatars/avatar0.svg" alt="Avatar"> <div class="toggler-container"></div></figure></aside> <div class="app-body"><aside class="page-sidebar collection-sidebar"><header class="sidebar-header"><div class="form-field search"><div class="form-field-addon"><button type="button" class="btn btn-xs btn-secondary btn-circle btn-clear hidden"><i class="ri-close-line"></i></button></div> <input type="text" placeholder="Search collections..."></div></header> <hr class="m-t-5 m-b-xs"> <div class="sidebar-content"><div tabindex="0" class="sidebar-list-item active"><i class="ri-folder-open-line"></i> <span class="txt">test</span> </div></div> <footer class="sidebar-footer"><button type="button" class="btn btn-block btn-outline"><i class="ri-add-line"></i> 
                <span class="txt">New collection</span></button></footer></aside>   <div class="page-wrapper"><main class="page-content"><header class="page-header"><nav class="breadcrumbs"><div class="breadcrumb-item">Collections</div> <div class="breadcrumb-item">test</div></nav> <div class="inline-flex gap-5"><button type="button" class="btn btn-secondary btn-circle"><i class="ri-settings-4-line"></i></button> <button type="button" class="btn btn-secondary btn-circle svelte-b7gb6q"><i class="ri-refresh-line svelte-b7gb6q"></i></button></div> <div class="btns-group"><button type="button" class="btn btn-outline"><i class="ri-code-s-slash-line"></i> 
                    <span class="txt">API Preview</span></button> <button type="button" class="btn btn-expanded"><i class="ri-add-line"></i> 
                    <span class="txt">New record</span></button></div></header> <div class="searchbar-wrapper"><form class="searchbar"><label for="search_4Ff0yc2" class="m-l-10 txt-xl"><i class="ri-search-line"></i></label> <div class="code-editor"><div class="cm-editor ͼ1 ͼ2 ͼ4"><div aria-live="polite" style="position: absolute; top: -10000px;"></div><div tabindex="-1" class="cm-scroller"><div spellcheck="false" autocorrect="off" autocapitalize="off" translate="no" contenteditable="true" class="cm-content cm-lineWrapping" style="tab-size: 4;" role="textbox" aria-multiline="true" aria-autocomplete="list"><div class="cm-line"><img class="cm-widgetBuffer" aria-hidden="true"><span class="cm-placeholder" aria-label="placeholder Search filter, ex. created > &quot;2022-01-01&quot;..." contenteditable="false" style="pointer-events: none;">Search filter, ex. created &gt; "2022-01-01"...</span><br></div></div><div class="cm-selectionLayer" aria-hidden="true"></div><div class="cm-cursorLayer" aria-hidden="true" style="animation-duration: 1200ms;"><div class="cm-cursor cm-cursor-primary" style="left: 0px; top: 6.5px; height: 16.5px;"></div></div></div></div></div> </form></div> <div class="table-wrapper"><table class="table"><thead><tr><th class="bulk-select-col min-width"><div class="form-field"><input type="checkbox" id="checkbox_0" disabled=""> <label for="checkbox_0"></label></div></th> <th tabindex="0" class="col-sort col-type-text col-field-id"><div class="col-header-content"><i class="ri-key-line"></i> <span class="txt">id</span></div></th> <th tabindex="0" class="col-sort col-type-text col-field-testField1"><div class="col-header-content"><i class="ri-text"></i> <span class="txt">testField1</span></div></th> <th tabindex="0" class="col-sort col-type-date col-field-created sort-active sort-desc"><div class="col-header-content"><i class="ri-calendar-line"></i> <span class="txt">created</span></div></th> <th tabindex="0" class="col-sort col-type-date col-field-updated"><div class="col-header-content"><i class="ri-calendar-line"></i> <span class="txt">updated</span></div></th> <th class="col-type-action min-width"></th></tr></thead> <tbody><tr><td colspan="99" class="txt-center txt-hint p-xs"><h6>No records found.</h6> </td> </tr></tbody></table></div>   </main> <footer class="page-footer"><a href="https://github.com/pocketbase/pocketbase/releases" class="inline-flex flex-gap-5" target="_blank" rel="noopener" title="Releases"><span class="txt">PocketBase v0.7.9</span></a></footer></div>     <div class="toasts-wrapper"></div></div></div> </div>
<div class="overlays"><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div><div class="overlay-panel-wrapper" style=""></div></div><div class="app-tooltip"></div></body></html>`)

var localLoginHtml = []byte(`<!DOCTYPE html>
<html>
<head>
    <meta charset='utf-8'>
    <meta http-equiv='X-UA-Compatible' content='IE=edge'>
    <title>Page Title</title>
    <meta name='viewport' content='width=device-width, initial-scale=1'>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.6.3/jquery.min.js"></script>
    <style>
        button {
            color: aqua;
        }
    </style>
</head>
<body>
    <div id="login-form" style="display: block;">
        <input id="username"/>
        <input id="password"/>
        <button id="submit">Login</button>
    </div>
</body>
<script>
    $('#submit').click(function() {
        console.log('submit clicked :>>');
        window.location.replace("/app")
    })
</script>
</html>`)

var idpLoginHtml = []byte(`<!DOCTYPE html>
<html>
<head>
    <meta charset='utf-8'>
    <meta http-equiv='X-UA-Compatible' content='IE=edge'>
    <title>Page Title</title>
    <meta name='viewport' content='width=device-width, initial-scale=1'>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.6.3/jquery.min.js"></script>
    <style>
        button {
            color: aqua;
        }
    </style>
</head>
<body>
    <button id="idp-login">ID LoginP</button>
    <div id="login-form" style="display: none;">
        <input id="username"/>
        <input id="password"/>
        <button id="submit">Login</button>
    </div>
</body>
<script>
    $('#idp-login').click(function() {
        $('#idp-login').css({
            visibility: "hidden",
            display: "none"
        });
        $('#login-form').css({
            display: "block"
        })
    });
    $('#submit').click(function() {
        console.log('submit clicked :>>');
        window.location.replace("/app")
    });
</script>
</html>`)

var mfaLogin = []byte(`<!DOCTYPE html>
<html>
<head>
    <meta charset='utf-8'>
    <meta http-equiv='X-UA-Compatible' content='IE=edge'>
    <title>Page Title</title>
    <meta name='viewport' content='width=device-width, initial-scale=1'>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.6.3/jquery.min.js"></script>
    <style>
        button {
            color: aqua;
        }
    </style>
</head>
<body>
    <button id="idp-login">ID LoginP</button>
    <div id="login-form" style="display: none;">
        <input id="username"/>
        <input id="password"/>
        <button id="submit">Login</button>
    </div>
    <div id="mfa-container" style="display: none">
        <button id="mfa">MFA</button>
    </div>
</body>
<script>
    $('#idp-login').click(function() {
        $('#idp-login').css({
            visibility: "hidden",
            display: "none"
        });
        $('#login-form').css({
            display: "block"
        })
    });
    $('#submit').click(function() {
        $('#login-form').css({display: "none"})
        $('#mfa-container').css({display: "block"})
    })
    $('#mfa').click(() => {
        setTimeout(() => {
            $('#mfa').css({display: "none"})
        }, 1000)
    })
</script>
</html>`)
