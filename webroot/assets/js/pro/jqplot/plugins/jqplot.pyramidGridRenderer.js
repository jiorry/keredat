/**
 * jqPlot
 * Pure JavaScript plotting plugin using jQuery
 *
 * Version: 1.0.8
 * Revision: 1250
 *
 * Copyright (c) 2009-2013 Chris Leonello
 * jqPlot is currently available for use in all personal or commercial projects 
 * under both the MIT (http://www.opensource.org/licenses/mit-license.php) and GPL 
 * version 2.0 (http://www.gnu.org/licenses/gpl-2.0.html) licenses. This means that you can 
 * choose the license that best suits your project and use it accordingly. 
 *
 * Although not required, the author would appreciate an email letting him 
 * know of any substantial use of jqPlot.  You can reach the author at: 
 * chris at jqplot dot com or see http://www.jqplot.com/info.php .
 *
 * If you are feeling kind and generous, consider supporting the project by
 * making a donation at: http://www.jqplot.com/donate.php .
 *
 * sprintf functions contained in jqplot.sprintf.js by Ash Searle:
 *
 *     version 2007.04.27
 *     author Ash Searle
 *     http://hexmen.com/blog/2007/03/printf-sprintf/
 *     http://hexmen.com/js/sprintf.js
 *     The author (Ash Searle) has placed this code in the public domain:
 *     "This code is unrestricted: you are free to use it however you like."
 * 
 */

(function(e){e.jqplot.PyramidGridRenderer=function(){e.jqplot.CanvasGridRenderer.call(this)},e.jqplot.PyramidGridRenderer.prototype=new e.jqplot.CanvasGridRenderer,e.jqplot.PyramidGridRenderer.prototype.constructor=e.jqplot.PyramidGridRenderer,e.jqplot.CanvasGridRenderer.prototype.init=function(t){this._ctx,this.plotBands={show:!1,color:"rgb(230, 219, 179)",axis:"y",start:null,interval:10},e.extend(!0,this,t);var n={lineJoin:"miter",lineCap:"round",fill:!1,isarc:!1,angle:this.shadowAngle,offset:this.shadowOffset,alpha:this.shadowAlpha,depth:this.shadowDepth,lineWidth:this.shadowWidth,closePath:!1,strokeStyle:this.shadowColor};this.renderer.shadowRenderer.init(n)},e.jqplot.PyramidGridRenderer.prototype.draw=function(){function L(n,r,i,s,o){t.save(),o=o||{};if(o.lineWidth==null||o.lineWidth!=0)e.extend(!0,t,o),t.beginPath(),t.moveTo(n,r),t.lineTo(i,s),t.stroke();t.restore()}this._ctx=this._elem.get(0).getContext("2d");var t=this._ctx,n=this._axes,r=n.xaxis.u2p,i=n.yMidAxis.u2p,s=n.xaxis.max/1e3,o=r(0),u=r(s),a=["xaxis","yaxis","x2axis","y2axis","yMidAxis"];t.save(),t.clearRect(0,0,this._plotDimensions.width,this._plotDimensions.height),t.fillStyle=this.backgroundColor||this.background,t.fillRect(this._left,this._top,this._width,this._height);if(this.plotBands.show){t.save();var f=this.plotBands;t.fillStyle=f.color;var l,c,h,p,d;f.axis.charAt(0)==="x"?n.xaxis.show&&(l=n.xaxis):f.axis.charAt(0)==="y"&&(n.yaxis.show?l=n.yaxis:n.y2axis.show?l=n.y2axis:n.yMidAxis.show&&(l=n.yMidAxis));if(l!==undefined){var v=f.start;v===null&&(v=l.min);for(var m=v;m<l.max;m+=2*f.interval)l.name.charAt(0)==="y"&&(c=this._left,m+f.interval<l.max?h=l.series_u2p(m+f.interval)+this._top:h=l.series_u2p(l.max)+this._top,p=this._right-this._left,d=l.series_u2p(v)-l.series_u2p(v+f.interval),t.fillRect(c,h,p,d))}t.restore()}t.save(),t.lineJoin="miter",t.lineCap="butt",t.lineWidth=this.gridLineWidth,t.strokeStyle=this.gridLineColor;var g,y,b,w;for(var m=5;m>0;m--){var E=a[m-1],l=n[E],S=l._ticks,x=S.length;if(l.show){if(l.drawBaseline){var T={};l.baselineWidth!==null&&(T.lineWidth=l.baselineWidth),l.baselineColor!==null&&(T.strokeStyle=l.baselineColor);switch(E){case"xaxis":n.yMidAxis.show?(L(this._left,this._bottom,o,this._bottom,T),L(u,this._bottom,this._right,this._bottom,T)):L(this._left,this._bottom,this._right,this._bottom,T);break;case"yaxis":L(this._left,this._bottom,this._left,this._top,T);break;case"yMidAxis":L(o,this._bottom,o,this._top,T),L(u,this._bottom,u,this._top,T);break;case"x2axis":n.yMidAxis.show?(L(this._left,this._top,o,this._top,T),L(u,this._top,this._right,this._top,T)):L(this._left,this._bottom,this._right,this._bottom,T);break;case"y2axis":L(this._right,this._bottom,this._right,this._top,T)}}for(var N=x;N>0;N--){var C=S[N-1];if(C.show){var k=Math.round(l.u2p(C.value))+.5;switch(E){case"xaxis":C.showGridline&&this.drawGridlines&&(!C.isMinorTick||l.showMinorTicks)&&L(k,this._top,k,this._bottom);if(C.showMark&&C.mark&&(!C.isMinorTick||l.showMinorTicks)){b=C.markSize,w=C.mark;var k=Math.round(l.u2p(C.value))+.5;switch(w){case"outside":g=this._bottom,y=this._bottom+b;break;case"inside":g=this._bottom-b,y=this._bottom;break;case"cross":g=this._bottom-b,y=this._bottom+b;break;default:g=this._bottom,y=this._bottom+b}this.shadow&&this.renderer.shadowRenderer.draw(t,[[k,g],[k,y]],{lineCap:"butt",lineWidth:this.gridLineWidth,offset:this.gridLineWidth*.75,depth:2,fill:!1,closePath:!1}),L(k,g,k,y)}break;case"yaxis":C.showGridline&&this.drawGridlines&&(!C.isMinorTick||l.showMinorTicks)&&L(this._right,k,this._left,k);if(C.showMark&&C.mark&&(!C.isMinorTick||l.showMinorTicks)){b=C.markSize,w=C.mark;var k=Math.round(l.u2p(C.value))+.5;switch(w){case"outside":g=this._left-b,y=this._left;break;case"inside":g=this._left,y=this._left+b;break;case"cross":g=this._left-b,y=this._left+b;break;default:g=this._left-b,y=this._left}this.shadow&&this.renderer.shadowRenderer.draw(t,[[g,k],[y,k]],{lineCap:"butt",lineWidth:this.gridLineWidth*1.5,offset:this.gridLineWidth*.75,fill:!1,closePath:!1}),L(g,k,y,k,{strokeStyle:l.borderColor})}break;case"yMidAxis":C.showGridline&&this.drawGridlines&&(!C.isMinorTick||l.showMinorTicks)&&(L(this._left,k,o,k),L(u,k,this._right,k));if(C.showMark&&C.mark&&(!C.isMinorTick||l.showMinorTicks)){b=C.markSize,w=C.mark;var k=Math.round(l.u2p(C.value))+.5;g=o,y=o+b,this.shadow&&this.renderer.shadowRenderer.draw(t,[[g,k],[y,k]],{lineCap:"butt",lineWidth:this.gridLineWidth*1.5,offset:this.gridLineWidth*.75,fill:!1,closePath:!1}),L(g,k,y,k,{strokeStyle:l.borderColor}),g=u-b,y=u,this.shadow&&this.renderer.shadowRenderer.draw(t,[[g,k],[y,k]],{lineCap:"butt",lineWidth:this.gridLineWidth*1.5,offset:this.gridLineWidth*.75,fill:!1,closePath:!1}),L(g,k,y,k,{strokeStyle:l.borderColor})}break;case"x2axis":C.showGridline&&this.drawGridlines&&(!C.isMinorTick||l.showMinorTicks)&&L(k,this._bottom,k,this._top);if(C.showMark&&C.mark&&(!C.isMinorTick||l.showMinorTicks)){b=C.markSize,w=C.mark;var k=Math.round(l.u2p(C.value))+.5;switch(w){case"outside":g=this._top-b,y=this._top;break;case"inside":g=this._top,y=this._top+b;break;case"cross":g=this._top-b,y=this._top+b;break;default:g=this._top-b,y=this._top}this.shadow&&this.renderer.shadowRenderer.draw(t,[[k,g],[k,y]],{lineCap:"butt",lineWidth:this.gridLineWidth,offset:this.gridLineWidth*.75,depth:2,fill:!1,closePath:!1}),L(k,g,k,y)}break;case"y2axis":C.showGridline&&this.drawGridlines&&(!C.isMinorTick||l.showMinorTicks)&&L(this._left,k,this._right,k);if(C.showMark&&C.mark&&(!C.isMinorTick||l.showMinorTicks)){b=C.markSize,w=C.mark;var k=Math.round(l.u2p(C.value))+.5;switch(w){case"outside":g=this._right,y=this._right+b;break;case"inside":g=this._right-b,y=this._right;break;case"cross":g=this._right-b,y=this._right+b;break;default:g=this._right,y=this._right+b}this.shadow&&this.renderer.shadowRenderer.draw(t,[[g,k],[y,k]],{lineCap:"butt",lineWidth:this.gridLineWidth*1.5,offset:this.gridLineWidth*.75,fill:!1,closePath:!1}),L(g,k,y,k,{strokeStyle:l.borderColor})}break;default:}}}C=null}l=null,S=null}t.restore();if(this.shadow)if(n.yMidAxis.show){var A=[[this._left,this._bottom],[o,this._bottom]];this.renderer.shadowRenderer.draw(t,A);var A=[[u,this._bottom],[this._right,this._bottom],[this._right,this._top]];this.renderer.shadowRenderer.draw(t,A);var A=[[o,this._bottom],[o,this._top]];this.renderer.shadowRenderer.draw(t,A)}else{var A=[[this._left,this._bottom],[this._right,this._bottom],[this._right,this._top]];this.renderer.shadowRenderer.draw(t,A)}this.borderWidth!=0&&this.drawBorder&&(n.yMidAxis.show?(L(this._left,this._top,o,this._top,{lineCap:"round",strokeStyle:n.x2axis.borderColor,lineWidth:n.x2axis.borderWidth}),L(u,this._top,this._right,this._top,{lineCap:"round",strokeStyle:n.x2axis.borderColor,lineWidth:n.x2axis.borderWidth}),L(this._right,this._top,this._right,this._bottom,{lineCap:"round",strokeStyle:n.y2axis.borderColor,lineWidth:n.y2axis.borderWidth}),L(this._right,this._bottom,u,this._bottom,{lineCap:"round",strokeStyle:n.xaxis.borderColor,lineWidth:n.xaxis.borderWidth}),L(o,this._bottom,this._left,this._bottom,{lineCap:"round",strokeStyle:n.xaxis.borderColor,lineWidth:n.xaxis.borderWidth}),L(this._left,this._bottom,this._left,this._top,{lineCap:"round",strokeStyle:n.yaxis.borderColor,lineWidth:n.yaxis.borderWidth}),L(o,this._bottom,o,this._top,{lineCap:"round",strokeStyle:n.yaxis.borderColor,lineWidth:n.yaxis.borderWidth}),L(u,this._bottom,u,this._top,{lineCap:"round",strokeStyle:n.yaxis.borderColor,lineWidth:n.yaxis.borderWidth})):(L(this._left,this._top,this._right,this._top,{lineCap:"round",strokeStyle:n.x2axis.borderColor,lineWidth:n.x2axis.borderWidth}),L(this._right,this._top,this._right,this._bottom,{lineCap:"round",strokeStyle:n.y2axis.borderColor,lineWidth:n.y2axis.borderWidth}),L(this._right,this._bottom,this._left,this._bottom,{lineCap:"round",strokeStyle:n.xaxis.borderColor,lineWidth:n.xaxis.borderWidth}),L(this._left,this._bottom,this._left,this._top,{lineCap:"round",strokeStyle:n.yaxis.borderColor,lineWidth:n.yaxis.borderWidth}))),t.restore(),t=null,n=null}})(jQuery);