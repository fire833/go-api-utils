/*
*	Copyright (C) 2025 Kendall Tauser
*
*	This program is free software; you can redistribute it and/or modify
*	it under the terms of the GNU General Public License as published by
*	the Free Software Foundation; either version 2 of the License, or
*	(at your option) any later version.
*
*	This program is distributed in the hope that it will be useful,
*	but WITHOUT ANY WARRANTY; without even the implied warranty of
*	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*	GNU General Public License for more details.
*
*	You should have received a copy of the GNU General Public License along
*	with this program; if not, write to the Free Software Foundation, Inc.,
*	51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

package object

import "google.golang.org/protobuf/reflect/protoreflect"

// Object is the primary interface for data objects that are managed and served by applications. Each Object, assuming
// it has the correct comment tags, will have an ObjectManager that is generated and imported into the application
// to allow for HTTP and Swagger support for the object. The manager will also have metrics for that object that can
// be updated and collected at runtime, including counts for each operation completed by that manager (ie number of GETs,
// PUTs, etc.)
type Object interface {
	// ProtoMessage allows Objects to be marshalled to protobuf automatically. This allows for an additional
	// MIME type to be supported by the API for free, essentially.
	protoreflect.ProtoMessage
}

type ObjectList interface {
	protoreflect.ProtoMessage

	GetItems() []Object
}
