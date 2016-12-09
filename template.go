package main

const templStreet1 string = `<svg id="Layer_1" data-name="Layer 1" 
    xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1247.24 609.45">
    <defs>
        <style>.cls-1{fill:#262741;}.cls-2{font-size:220px;font-family:ProbaNav2-SemiBold, Proba Nav2;font-weight:700;letter-spacing:0em;}.cls-10,.cls-14,.cls-2{fill:#fff;}.cls-3{letter-spacing:-0.01em;}.cls-4{letter-spacing:0em;}.cls-5{letter-spacing:0em;}.cls-6{letter-spacing:-0.01em;}.cls-7{letter-spacing:0em;}.cls-8{letter-spacing:-0.01em;}.cls-9{letter-spacing:0.01em;}.cls-10,.cls-14{font-size:90px;font-family:ProbaNav2-Regular, Proba Nav2;}.cls-10{letter-spacing:0em;}.cls-11{letter-spacing:0.01em;}.cls-12,.cls-14{letter-spacing:0em;}.cls-13{letter-spacing:0.01em;}.cls-15{letter-spacing:-0.01em;}.cls-16{fill:none;stroke:#fff;stroke-miterlimit:10;stroke-width:4px;}</style>
    </defs>
    <title>Artboard 2</title>
    <rect class="cls-1" x="318.9" y="-318.9" width="609.45" height="1247.24" rx="14.17" ry="14.17" transform="translate(928.35 -318.9) rotate(90)"/>
    <text class="cls-2" transform="translate(98.42 340.64)">
        {{ .StreetNameUA }}
    </text>
    <text class="cls-10" transform="translate(107.32 537.78)">
       {{ .StreetNameEng }} {{ .StreetType }}
    </text>
    <text class="cls-14" transform="translate(107.3 118.27)">
        {{ .StreetTypeUA }}
    </text>
    <line class="cls-16" x1="110.09" y1="431.55" x2="1133.86" y2="431.55"/>
</svg>`

const templStreet2 string = `<svg id="Layer_1" data-name="Layer 1" 
    xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1785.83 609.45">
    <defs>
        <style>.cls-1{fill:#262741;}.cls-2{font-size:220px;font-family:ProbaNav2-SemiBold, Proba Nav2;font-weight:700;letter-spacing:-0.02em;}.cls-12,.cls-18,.cls-2{fill:#fff;}.cls-3{letter-spacing:0em;}.cls-4{letter-spacing:0em;}.cls-5{letter-spacing:-0.02em;}.cls-6{letter-spacing:-0.01em;}.cls-7{letter-spacing:-0.01em;}.cls-8{letter-spacing:0em;}.cls-9{letter-spacing:-0.02em;}.cls-10{letter-spacing:-0.01em;}.cls-11{letter-spacing:0.01em;}.cls-12,.cls-18{font-size:90px;font-family:ProbaNav2-Regular, Proba Nav2;}.cls-12{letter-spacing:0.01em;}.cls-13{letter-spacing:0.01em;}.cls-14{letter-spacing:0.01em;}.cls-15{letter-spacing:-0.02em;}.cls-16{letter-spacing:-0.01em;}.cls-17{letter-spacing:0.01em;}.cls-18{letter-spacing:0em;}.cls-19{letter-spacing:-0.01em;}.cls-20{fill:none;stroke:#fff;stroke-miterlimit:10;stroke-width:4px;}</style>
    </defs>
    <title>Artboard 3</title>
    <path class="cls-1" d="M588.19-588.19h595.28A14.17,14.17,0,0,1,1197.64-574V1183.46a14.17,14.17,0,0,1-14.17,14.17H602.36a14.17,14.17,0,0,1-14.17-14.17V-588.19A0,0,0,0,1,588.19-588.19Z" transform="translate(1197.64 -588.19) rotate(90)"/>
    <text class="cls-2" transform="translate(101.72 340.64)">
        {{ .StreetNameUA }}
    </text>
    <text class="cls-12" transform="translate(110.61 537.78)">
        {{ .StreetNameEng }} {{ .StreetType }}
    </text>
    <text class="cls-18" transform="translate(110.59 118.27)">
        {{ .StreetTypeUA }}
    </text>
    <line class="cls-20" x1="113.39" y1="431.55" x2="1672.44" y2="431.55"/>
</svg>`

const templStreet3 string = `<svg id="Layer_1" data-name="Layer 1" 
    xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2239.37 609.45">
    <defs>
        <style>.cls-1{fill:#262741;}.cls-2{font-size:220px;font-family:ProbaNav2-SemiBold, Proba Nav2;font-weight:700;letter-spacing:-0.01em;}.cls-17,.cls-2,.cls-21{fill:#fff;}.cls-3{letter-spacing:0em;}.cls-4{letter-spacing:-0.04em;}.cls-5{letter-spacing:-0.01em;}.cls-6{letter-spacing:0em;}.cls-7{letter-spacing:-0.01em;}.cls-8{letter-spacing:0.01em;}.cls-9{letter-spacing:-0.01em;}.cls-10{letter-spacing:0.01em;}.cls-11{letter-spacing:0em;}.cls-12{letter-spacing:-0.01em;}.cls-13{letter-spacing:0em;}.cls-14{letter-spacing:-0.01em;}.cls-15{letter-spacing:0em;}.cls-16{letter-spacing:-0.01em;}.cls-17,.cls-21{font-size:90px;font-family:ProbaNav2-Regular, Proba Nav2;}.cls-18{letter-spacing:0.01em;}.cls-19{letter-spacing:-0.01em;}.cls-20{letter-spacing:0em;}.cls-21{letter-spacing:0em;}.cls-22{letter-spacing:-0.01em;}.cls-23{letter-spacing:0.01em;}.cls-24{fill:none;stroke:#fff;stroke-miterlimit:10;stroke-width:4px;}</style>
    </defs>
    <title>Artboard 4</title>
    <path class="cls-1" d="M815-815h595.28a14.17,14.17,0,0,1,14.17,14.17v2211a14.17,14.17,0,0,1-14.17,14.17H829.13A14.17,14.17,0,0,1,815,1410.24V-815a0,0,0,0,1,0,0Z" transform="translate(1424.41 -814.96) rotate(90)"/>
    <text class="cls-2" transform="translate(98.88 340.64)">
        {{ .StreetNameUA }}
    </text>
    <text class="cls-17" transform="translate(107.77 537.78)">
        {{ .StreetNameEng }} {{ .StreetType }}
    </text>
    <text class="cls-21" transform="translate(107.76 118.27)">
        {{ .StreetTypeUA }}
    </text>
    <line class="cls-24" x1="110.55" y1="431.55" x2="2125.98" y2="431.55"/>
</svg>`

const templStreet4 string = `<svg id="Layer_1" data-name="Layer 1" 
    xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2806.3 609.45">
    <defs>
        <style>.cls-1{fill:#262741;}.cls-2{font-size:220px;font-family:ProbaNav2-SemiBold, Proba Nav2;font-weight:700;letter-spacing:-0.01em;}.cls-18,.cls-2,.cls-22{fill:#fff;}.cls-3{letter-spacing:-0.01em;}.cls-4{letter-spacing:-0.04em;}.cls-5{letter-spacing:0em;}.cls-6{letter-spacing:0em;}.cls-7{letter-spacing:-0.01em;}.cls-8{letter-spacing:0.01em;}.cls-9{letter-spacing:-0.01em;}.cls-10{letter-spacing:0em;}.cls-11{letter-spacing:-0.01em;}.cls-12{letter-spacing:0em;}.cls-13{letter-spacing:-0.02em;}.cls-14{letter-spacing:0em;}.cls-15{letter-spacing:-0.04em;}.cls-16{letter-spacing:-0.02em;}.cls-17{letter-spacing:0.01em;}.cls-18,.cls-22{font-size:90px;font-family:ProbaNav2-Regular, Proba Nav2;}.cls-18,.cls-24{letter-spacing:0.01em;}.cls-19{letter-spacing:-0.01em;}.cls-20{letter-spacing:0.01em;}.cls-21{letter-spacing:-0.01em;}.cls-22{letter-spacing:0em;}.cls-23{letter-spacing:-0.01em;}.cls-25{fill:none;stroke:#fff;stroke-miterlimit:10;stroke-width:4px;}</style>
    </defs>
    <title>Artboard 5</title>
    <path class="cls-1" d="M1098.43-1098.43H1693.7a14.17,14.17,0,0,1,14.17,14.17v2778a14.17,14.17,0,0,1-14.17,14.17H1112.6a14.17,14.17,0,0,1-14.17-14.17V-1098.43A0,0,0,0,1,1098.43-1098.43Z" transform="translate(1707.87 -1098.43) rotate(90)"/>
    <text class="cls-2" transform="translate(98.88 340.64)">
        {{ .StreetNameUA }}
    </text>
    <text class="cls-18" transform="translate(107.77 537.78)">
        {{ .StreetNameEng }} {{ .StreetType }}
    </text>
    <text class="cls-22" transform="translate(107.76 118.27)">
        {{ .StreetTypeUA }}
    </text>
    <line class="cls-25" x1="110.55" y1="431.55" x2="2693.42" y2="431.55"/>
</svg>`
