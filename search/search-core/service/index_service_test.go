package service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IndexService", func() {
    
    Describe("textSlicePos", func() {
        
        It("returns -1 when item is not found", func() {
            pos := textSlicePos([]string{"1", "2"}, "0")
            Expect(pos).To(Equal(-1))
        })
        
        It("returns -1 when slice is empty", func() {
            var slice []string
            pos := textSlicePos(slice, "0")
            Expect(pos).To(Equal(-1))
        })

        It("returns correct pos when text is in slice", func() {
            slice := []string{"1","2","3"}
            
            pos := textSlicePos(slice, "1")
            Expect(pos).To(Equal(0))

            pos = textSlicePos(slice, "2")
            Expect(pos).To(Equal(1))

            pos = textSlicePos(slice, "3")
            Expect(pos).To(Equal(2))
        })

    })

    Describe("textSliceRemovePos", func() {
        
        It("removes first element", func() {
            slice := textSliceRemovePos([]string{"1","2","3"}, 0)
            Expect(len(slice)).To(Equal(2))
            Expect(slice[0]).To(Equal("2"))
        })

        It("removes last element", func() {
            slice := textSliceRemovePos([]string{"1","2","3"}, 2)
            Expect(len(slice)).To(Equal(2))
            Expect(slice[len(slice)-1]).To(Equal("2"))
        })

        It("removes middle element", func() {
            slice := textSliceRemovePos([]string{"1","2","3"}, 1)
            Expect(len(slice)).To(Equal(2))
            Expect(slice[1]).To(Equal("3"))
        })


        It("does nothing if pos is outside boundaries", func() {
            slice := textSliceRemovePos([]string{"1","2","3"}, -1)
            Expect(len(slice)).To(Equal(3))
            slice = textSliceRemovePos([]string{"1","2","3"}, 3)
            Expect(len(slice)).To(Equal(3))
        })

    })


})
