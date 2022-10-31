#include "viewport.h"

#include "CMU462.h"

namespace CMU462 {

void ViewportImp::set_viewbox( float centerX, float centerY, float vspan ) {

  // Task 5 (part 2): 
  // Set svg coordinate to normalized device coordinate transformation. Your input
  // arguments are defined as normalized SVG canvas coordinates.
  this->centerX = centerX;
  this->centerY = centerY;
  this->vspan = vspan; 

  float x_low_dash = 0;
  float x_high_dash = 1;
  float y_low_dash = 0;
  float y_high_dash = 1;

  float x_high = centerX + vspan; // using vspan as offset
  float x_low = centerX - vspan;
  float y_high = centerY + vspan;
  float y_low = centerY - vspan;

  Matrix3x3 Tmat; // transformation matrix

  Tmat(0, 0) = (x_high_dash - x_low_dash) / (x_high - x_low);
  Tmat(0, 1) = 0;
  Tmat(0, 2) = ((x_low_dash * x_high) - (x_high_dash * x_low))/ (x_high - x_low);
  Tmat(1, 0) = 0;
  Tmat(1, 1) = (y_high_dash - y_low_dash) / (y_high - y_low);
  Tmat(1, 2) = ((y_low_dash * y_high) - (y_high_dash * y_low)) / (y_high - y_low);
  Tmat(2, 0) = 0;
  Tmat(2, 1) = 0;
  Tmat(2, 2) = 1;

  set_svg_2_norm(Tmat);
}

void ViewportImp::update_viewbox( float dx, float dy, float scale ) { 
  
  this->centerX -= dx;
  this->centerY -= dy;
  this->vspan *= scale;
  set_viewbox( centerX, centerY, vspan );
}

} // namespace CMU462
