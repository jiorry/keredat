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

(function(e){function t(t,n,r){r=r||{},r.axesDefaults=r.axesDefaults||{},r.grid=r.grid||{},r.legend=r.legend||{},r.seriesDefaults=r.seriesDefaults||{};var i=!1;if(r.seriesDefaults.renderer===e.jqplot.PyramidRenderer)i=!0;else if(r.series)for(var s=0;s<r.series.length;s++)r.series[s].renderer===e.jqplot.PyramidRenderer&&(i=!0);i&&(r.axesDefaults.renderer=e.jqplot.PyramidAxisRenderer,r.grid.renderer=e.jqplot.PyramidGridRenderer,r.seriesDefaults.pointLabels={show:!1})}function n(){this.plugins.pyramidRenderer&&this.plugins.pyramidRenderer.highlightCanvas&&(this.plugins.pyramidRenderer.highlightCanvas.resetCanvas(),this.plugins.pyramidRenderer.highlightCanvas=null),this.plugins.pyramidRenderer={highlightedSeriesIndex:null},this.plugins.pyramidRenderer.highlightCanvas=new e.jqplot.GenericCanvas,this.eventCanvas._elem.before(this.plugins.pyramidRenderer.highlightCanvas.createElement(this._gridPadding,"jqplot-pyramidRenderer-highlight-canvas",this._plotDimensions,this)),this.plugins.pyramidRenderer.highlightCanvas.setContext(),this.eventCanvas._elem.bind("mouseleave",{plot:this},function(e){i(e.data.plot)})}function r(e,t,n,r){var i=e.series[t],s=e.plugins.pyramidRenderer.highlightCanvas;s._ctx.clearRect(0,0,s._ctx.canvas.width,s._ctx.canvas.height),i._highlightedPoint=n,e.plugins.pyramidRenderer.highlightedSeriesIndex=t;var o={fillStyle:i.highlightColors[n],fillRect:!1};i.renderer.shapeRenderer.draw(s._ctx,r,o),i.synchronizeHighlight!==!1&&e.series.length>=i.synchronizeHighlight&&i.synchronizeHighlight!==t&&(i=e.series[i.synchronizeHighlight],o={fillStyle:i.highlightColors[n],fillRect:!1},i.renderer.shapeRenderer.draw(s._ctx,i._barPoints[n],o)),s=null}function i(e){var t=e.plugins.pyramidRenderer.highlightCanvas;t._ctx.clearRect(0,0,t._ctx.canvas.width,t._ctx.canvas.height);for(var n=0;n<e.series.length;n++)e.series[n]._highlightedPoint=null;e.plugins.pyramidRenderer.highlightedSeriesIndex=null,e.target.trigger("jqplotDataUnhighlight"),t=null}function s(e,t,n,s,o){if(s){var u=[s.seriesIndex,s.pointIndex,s.data],a=jQuery.Event("jqplotDataMouseOver");a.pageX=e.pageX,a.pageY=e.pageY,o.target.trigger(a,u);if(o.series[u[0]].highlightMouseOver&&(u[0]!=o.plugins.pyramidRenderer.highlightedSeriesIndex||u[1]!=o.series[u[0]]._highlightedPoint)){var f=jQuery.Event("jqplotDataHighlight");f.which=e.which,f.pageX=e.pageX,f.pageY=e.pageY,o.target.trigger(f,u),r(o,s.seriesIndex,s.pointIndex,s.points)}}else s==null&&i(o)}e.jqplot.PyramidAxisRenderer===undefined&&e.ajax({url:e.jqplot.pluginLocation+"jqplot.pyramidAxisRenderer.js",dataType:"script",async:!1}),e.jqplot.PyramidGridRenderer===undefined&&e.ajax({url:e.jqplot.pluginLocation+"jqplot.pyramidGridRenderer.js",dataType:"script",async:!1}),e.jqplot.PyramidRenderer=function(){e.jqplot.LineRenderer.call(this)},e.jqplot.PyramidRenderer.prototype=new e.jqplot.LineRenderer,e.jqplot.PyramidRenderer.prototype.constructor=e.jqplot.PyramidRenderer,e.jqplot.PyramidRenderer.prototype.init=function(t,r){t=t||{},this._type="pyramid",this.barPadding=10,this.barWidth=null,this.fill=!0,this.highlightMouseOver=!0,this.highlightMouseDown=!1,this.highlightColors=[],this.highlightThreshold=2,this.synchronizeHighlight=!1,this.offsetBars=!1,t.highlightMouseDown&&t.highlightMouseOver==null&&(t.highlightMouseOver=!1),this.side="right",e.extend(!0,this,t),this.side==="left"?this._highlightThreshold=[[-this.highlightThreshold,0],[-this.highlightThreshold,0],[0,0],[0,0]]:this._highlightThreshold=[[0,0],[0,0],[this.highlightThreshold,0],[this.highlightThreshold,0]],this.renderer.options=t,this._highlightedPoint=null,this._dataColors=[],this._barPoints=[],this.fillAxis="y",this._primaryAxis="_yaxis",this._xnudge=0;var i={lineJoin:"miter",lineCap:"butt",fill:this.fill,fillRect:this.fill,isarc:!1,strokeStyle:this.color,fillStyle:this.color,closePath:this.fill,lineWidth:this.lineWidth};this.renderer.shapeRenderer.init(i);var o=t.shadowOffset;o==null&&(this.lineWidth>2.5?o=1.25*(1+(Math.atan(this.lineWidth/2.5)/.785398163-1)*.6):o=1.25*Math.atan(this.lineWidth/2.5)/.785398163);var u={lineJoin:"miter",lineCap:"butt",fill:this.fill,fillRect:this.fill,isarc:!1,angle:this.shadowAngle,offset:o,alpha:this.shadowAlpha,depth:this.shadowDepth,closePath:this.fill,lineWidth:this.lineWidth};this.renderer.shadowRenderer.init(u),r.postDrawHooks.addOnce(n),r.eventListenerHooks.addOnce("jqplotMouseMove",s);if(this.side==="left")for(var a=0,f=this.data.length;a<f;a++)this.data[a][1]=-Math.abs(this.data[a][1])},e.jqplot.PyramidRenderer.prototype.setGridData=function(e){var t=this._xaxis.series_u2p,n=this._yaxis.series_u2p,r=this._plotData,i=this._prevPlotData;this.gridData=[],this._prevGridData=[];var s=r.length,o=!1,u;for(u=0;u<s;u++)r[u][1]<0&&(this.side="left");this._yaxis.name==="yMidAxis"&&this.side==="right"&&(this._xnudge=this._xaxis.max/2e3,o=!0);for(u=0;u<s;u++)r[u][0]!=null&&r[u][1]!=null?this.gridData.push([t(r[u][1]),n(r[u][0])]):r[u][0]==null?this.gridData.push([t(r[u][1]),null]):r[u][1]==null&&this.gridData.push(null,[n(r[u][0])]),r[u][1]===0&&o&&(this.gridData[u][0]=t(this._xnudge))},e.jqplot.PyramidRenderer.prototype.makeGridData=function(e,t){var n=this._xaxis.series_u2p,r=this._yaxis.series_u2p,i=[],s=e.length,o=!1,u;for(u=0;u<s;u++)e[u][1]<0&&(this.side="left");this._yaxis.name==="yMidAxis"&&this.side==="right"&&(this._xnudge=this._xaxis.max/2e3,o=!0);for(u=0;u<s;u++)e[u][0]!=null&&e[u][1]!=null?i.push([n(e[u][1]),r(e[u][0])]):e[u][0]==null?i.push([n(e[u][1]),null]):e[u][1]==null&&i.push([null,r(e[u][0])]),e[u][1]===0&&o&&(i[u][0]=n(this._xnudge));return i},e.jqplot.PyramidRenderer.prototype.setBarWidth=function(){var e,t=0,n=0,r=this[this._primaryAxis],i,s,o;t=r.max-r.min;var u=r.numberTicks,a=(u-1)/2,f=this.barPadding===0?1:0;r.name=="xaxis"||r.name=="x2axis"?this.barWidth=(r._offsets.max-r._offsets.min)/t-this.barPadding+f:this.fill?this.barWidth=(r._offsets.min-r._offsets.max)/t-this.barPadding+f:this.barWidth=(r._offsets.min-r._offsets.max)/t},e.jqplot.PyramidRenderer.prototype.draw=function(t,n,r){var i,s=e.extend({},r),o=s.shadow!=undefined?s.shadow:this.shadow,u=s.showLine!=undefined?s.showLine:this.showLine,a=s.fill!=undefined?s.fill:this.fill,f=this._xaxis.series_u2p,l=this._yaxis.series_u2p,c,h;this._dataColors=[],this._barPoints=[],this.renderer.options.barWidth==null&&this.renderer.setBarWidth.call(this);var p=[],d,v;if(u){var m=new e.jqplot.ColorGenerator(this.negativeSeriesColors),g=new e.jqplot.ColorGenerator(this.seriesColors),y=m.get(this.index);this.useNegativeColors||(y=s.fillStyle);var b=s.fillStyle,w,E=this._xaxis.series_u2p(this._xnudge),S=this._yaxis.series_u2p(this._yaxis.min),x=this._yaxis.series_u2p(this._yaxis.max),T=this.barWidth,N=T/2,p=[],C=this.offsetBars?N:0;for(var i=0,k=n.length;i<k;i++){if(this.data[i][0]==null)continue;w=n[i][1],this._plotData[i][1]<0?this.varyBarColor&&!this._stack&&(this.useNegativeColors?s.fillStyle=m.next():s.fillStyle=g.next()):this.varyBarColor&&!this._stack?s.fillStyle=g.next():s.fillStyle=b;if(this.fill){this._plotData[i][1]>=0?(d=n[i][0]-E,v=this.barWidth,p=[E,w-N-C,d,v]):(d=E-n[i][0],v=this.barWidth,p=[n[i][0],w-N-C,d,v]),this._barPoints.push([[p[0],p[1]+v],[p[0],p[1]],[p[0]+d,p[1]],[p[0]+d,p[1]+v]]),o&&this.renderer.shadowRenderer.draw(t,p);var L=s.fillStyle||this.color;this._dataColors.push(L),this.renderer.shapeRenderer.draw(t,p,s)}else if(i===0)p=[[E,S],[n[i][0],S],[n[i][0],n[i][1]-N-C]];else if(i<k-1)p=p.concat([[n[i-1][0],n[i-1][1]-N-C],[n[i][0],n[i][1]+N-C],[n[i][0],n[i][1]-N-C]]);else{p=p.concat([[n[i-1][0],n[i-1][1]-N-C],[n[i][0],n[i][1]+N-C],[n[i][0],x],[E,x]]),o&&this.renderer.shadowRenderer.draw(t,p);var L=s.fillStyle||this.color;this._dataColors.push(L),this.renderer.shapeRenderer.draw(t,p,s)}}}if(this.highlightColors.length==0)this.highlightColors=e.jqplot.computeHighlightColors(this._dataColors);else if(typeof this.highlightColors=="string"){this.highlightColors=[];for(var i=0;i<this._dataColors.length;i++)this.highlightColors.push(this.highlightColors)}},e.jqplot.preInitHooks.push(t)})(jQuery);