//-----------------------------------------------------------------------------
/*

Common 3D shapes.

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------

// Counter Bored Hole
func CounterBored_Hole3D(
	l float64, // total length
	r float64, // hole radius
	cb_r float64, // counter bore radius
	cb_d float64, // counter bore depth
) SDF3 {
	s0 := Cylinder3D(l, r, 0)
	s1 := Cylinder3D(cb_d, cb_r, 0)
	s1 = Transform3D(s1, Translate3d(V3{0, 0, (l - cb_d) / 2}))
	return Union3D(s0, s1)
}

// Chamfered Hole (45 degrees)
func Chamfered_Hole3D(
	l float64, // total length
	r float64, // hole radius
	ch_r float64, // chamfer radius
) SDF3 {
	s0 := Cylinder3D(l, r, 0)
	s1 := Cone3D(ch_r, r, r+ch_r, 0)
	s1 = Transform3D(s1, Translate3d(V3{0, 0, (l - ch_r) / 2}))
	return Union3D(s0, s1)
}

// Countersunk Hole (45 degrees)
func CounterSunk_Hole3D(
	l float64, // total length
	r float64, // hole radius
) SDF3 {
	return Chamfered_Hole3D(l, r, r)
}

//-----------------------------------------------------------------------------

// Return a rounded hex head for a nut or bolt.
func HexHead3D(
	r float64, // radius
	h float64, // height
	round string, // (t)top, (b)bottom, (tb)top/bottom
) SDF3 {
	// basic hex body
	corner_round := r * 0.08
	hex_2d := Polygon2D(Nagon(6, r-corner_round))
	hex_2d = Offset2D(hex_2d, corner_round)
	hex_3d := Extrude3D(hex_2d, h)
	// round out the top and/or bottom as required
	if round != "" {
		top_round := r * 1.6
		d := r * math.Cos(DtoR(30))
		sphere_3d := Sphere3D(top_round)
		z_ofs := math.Sqrt(top_round*top_round-d*d) - h/2
		if round == "t" || round == "tb" {
			hex_3d = Intersect3D(hex_3d, Transform3D(sphere_3d, Translate3d(V3{0, 0, -z_ofs})))
		}
		if round == "b" || round == "tb" {
			hex_3d = Intersect3D(hex_3d, Transform3D(sphere_3d, Translate3d(V3{0, 0, z_ofs})))
		}
	}
	return hex_3d
}

// Return a cylindrical knurled head.
func KnurledHead3D(
	r float64, // radius
	h float64, // height
	pitch float64, // knurl pitch
) SDF3 {
	theta := DtoR(45)
	cylinder_round := r * 0.05
	knurl_h := pitch * math.Floor((h-cylinder_round)/pitch)
	knurl_3d := Knurl3D(knurl_h, r, pitch, pitch*0.3, theta)
	return Union3D(Cylinder3D(h, r, cylinder_round), knurl_3d)
}

//-----------------------------------------------------------------------------

// Return a 2D knurl profile.
func KnurlProfile(
	radius float64, // radius of knurled cylinder
	pitch float64, // pitch of the knurl
	height float64, // height of the knurl
) SDF2 {
	knurl := NewPolygon()
	knurl.Add(pitch/2, 0)
	knurl.Add(pitch/2, radius)
	knurl.Add(0, radius+height)
	knurl.Add(-pitch/2, radius)
	knurl.Add(-pitch/2, 0)
	//knurl.Render("knurl.dxf")
	return Polygon2D(knurl.Vertices())
}

// Return a knurled cylinder.
func Knurl3D(
	length float64, // length of cylinder
	radius float64, // radius of cylinder
	pitch float64, // knurl pitch
	height float64, // knurl height
	theta float64, // knurl helix angle
) SDF3 {
	// A knurl is the the intersection of left and right hand
	// multistart "threads". Work out the number of starts using
	// the desired helix angle.
	n := int(TAU * radius * math.Tan(theta) / pitch)
	// build the knurl profile.
	knurl_2d := KnurlProfile(radius, pitch, height)
	// create the left/right hand spirals
	knurl0_3d := Screw3D(knurl_2d, length, pitch, n)
	knurl1_3d := Screw3D(knurl_2d, length, pitch, -n)
	return Intersect3D(knurl0_3d, knurl1_3d)
}

//-----------------------------------------------------------------------------

// Return a washer.
func Washer3D(
	t float64, // thickness
	r_inner float64, // inner radius
	r_outer float64, // outer radius
) SDF3 {
	if t <= 0 {
		panic("t <= 0")
	}
	if r_inner >= r_outer {
		panic("r_inner >= r_outer")
	}
	return Difference3D(Cylinder3D(t, r_outer, 0), Cylinder3D(t, r_inner, 0))
}

//-----------------------------------------------------------------------------
// Board standoffs

type Standoff_Parms struct {
	Pillar_height float64
	Pillar_radius float64
	Hole_depth    float64
	Hole_radius   float64
	Number_webs   int
	Web_height    float64
	Web_radius    float64
	Web_width     float64
}

// single web
func pillar_web(p *Standoff_Parms) SDF3 {
	w := NewPolygon()
	w.Add(0, 0)
	w.Add(p.Web_radius, 0)
	w.Add(0, p.Web_height)
	s := Extrude3D(Polygon2D(w.Vertices()), p.Web_width)
	m := Translate3d(V3{0, 0, -0.5 * p.Pillar_height}).Mul(RotateX(DtoR(90.0)))
	return Transform3D(s, m)
}

// multiple webs
func pillar_webs(p *Standoff_Parms) SDF3 {
	if p.Number_webs == 0 {
		return nil
	}
	return RotateCopy3D(pillar_web(p), p.Number_webs)
}

// pillar
func pillar(p *Standoff_Parms) SDF3 {
	return Cylinder3D(p.Pillar_height, p.Pillar_radius, 0)
}

// pillar hole
func pillar_hole(p *Standoff_Parms) SDF3 {
	if p.Hole_radius == 0.0 || p.Hole_depth == 0.0 {
		return nil
	}
	s := Cylinder3D(p.Hole_depth, p.Hole_radius, 0)
	z_ofs := 0.5 * (p.Pillar_height - p.Hole_depth)
	return Transform3D(s, Translate3d(V3{0, 0, z_ofs}))
}

func Standoff3D(p *Standoff_Parms) SDF3 {
	s0 := Difference3D(Union3D(pillar(p), pillar_webs(p)), pillar_hole(p))
	if p.Number_webs != 0 {
		// Cut off any part of the webs that protrude from the top of the pillar
		s1 := Cylinder3D(p.Pillar_height, 2.0*p.Web_radius, 0)
		return Intersect3D(s0, s1)
	}
	return s0
}

//-----------------------------------------------------------------------------
// finger button

type FingerButtonParms struct {
	Size      V2      // button size
	Gap       float64 // gap between finger and body
	Length    float64 // length of the finger
	Thickness float64 // thickness of the button
}

func FingerButton3D(k *FingerButtonParms) SDF2 {

	// rounding radius based on button size
	r0 := 0.4 * Min(k.Size.X, k.Size.Y)
	r1 := r0 + k.Gap

	// finger width based on button Size.X
	fw := 0.7 * k.Size.X
	// finger y offset
	f_ofs := 0.5 * (k.Length + k.Size.Y)

	// 2d button (0,0)
	b := Box2D(k.Size, r0)

	// 2d button surround (0,0)
	bx := Box2D(k.Size.AddScalar(2.0*k.Gap), r1)

	// 2d finger
	f := Box2D(V2{fw, k.Length}, 0)
	f = Transform2D(f, Translate2d(V2{0, f_ofs}))

	// 2d finger surround
	fx := Box2D(V2{fw + (2.0 * k.Gap), k.Length}, r0)
	fx = Transform2D(fx, Translate2d(V2{0, f_ofs}))

	s := Difference2D(Union2D(bx, fx), Union2D(b, f))
	return s

	//return nil
}

//-----------------------------------------------------------------------------
