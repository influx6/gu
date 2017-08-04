package css

const (
	base = `
.line-height-print, .line-height-phone {
	line-height: 1.25;
}

.line-height-deskop, .line-height-desktop-large {
	line-height: 1.3756;
}

.line-height-tablet {
	line-height: 1.3756;
}


.letter-spacing-1 {
	letter-spacing: 1px;
}

.letter-spacing-2 {
	letter-spacing: -0.5px;
}

.letter-spacing-3 {
	letter-spacing: -1px;
}

.letter-spacing-4 {
	letter-spacing: -2px;
}

.focus-shadow {
  box-shadow: 0 0 8px rgba(0,0,0,.18),0 8px 16px rgba(0,0,0,.36);
}

.shadow-2dp {
	box-shadow: 0 2px 2px 0 rgba(0, 0, 0, 0.14),
              0 3px 1px -2px rgba(0, 0, 0, 0.2),
              0 1px 5px 0 rgba(0, 0, 0, 0.12);
}

.shadow-3dp {
  box-shadow: 0 3px 4px 0 rgba(0, 0, 0, 0.14),
              0 3px 3px -2px rgba(0, 0, 0, 0.2),
              0 1px 8px 0 rgba(0, 0, 0, 0.12);
}

.shadow-4dp {
  box-shadow: 0 4px 5px 0 rgba(0, 0, 0, 0.14),
              0 1px 10px 0 rgba(0, 0, 0, 0.12),
              0 2px 4px -1px rgba(0, 0, 0, 0.2);
}

.shadow-6dp {
  box-shadow: 0 6px 10px 0 rgba(0, 0, 0, 0.14),
              0 1px 18px 0 rgba(0, 0, 0, 0.12),
              0 3px 5px -1px rgba(0, 0, 0, 0.2);
}

.shadow-8dp {
  box-shadow: 0 8px 10px 1px rgba(0, 0, 0, 0.14),
              0 3px 14px 2px rgba(0, 0, 0, 0.12),
              0 5px 5px -3px rgba(0, 0, 0, 0.2);
}

.shadow-16dp {
  box-shadow: 0 16px 24px 2px rgba(0, 0, 0, 0.14),
              0  6px 30px 5px rgba(0, 0, 0, 0.12),
              0  8px 10px -5px rgba(0, 0, 0, 0.2);
}

.shadow-24dp {
  box-shadow: 0  9px 46px  8px rgba(0, 0, 0, 0.14),
              0 11px 15px -7px rgba(0, 0, 0, 0.12),
              0 24px 38px  3px rgba(0, 0, 0, 0.2);
}

.shadow {
	box-shadow: 0px 13px 20px 2px rgba(0, 0, 0, 0.45);
}

.shadow-dropdown {
	box-shadow: 0px 9px 30px 2px rgba(0, 0, 0, 0.51);
}

.shadow-hover {
	box-shadow: 0px 13px 30px 5px rgba(0, 0, 0, 0.58);
}

.shadow-elevated {
	box-shadow: 0px 20px 40px 4px rgba(0, 0, 0, 0.51);
}

.wrap {
  text-wrap: wrap;
  white-space: -moz-pre-wrap;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.sizing {
  box-sizing: border-box;
  -webkit-box-sizing: border-box;
  -o-box-sizing: border-box;
  -moz-box-sizing: border-box;
}

.clear {
	content: " ";
	clear: both;
	display: block;
	visibility: hidden;
	height: 0;
	font-size: 0;
}
	
`
)
